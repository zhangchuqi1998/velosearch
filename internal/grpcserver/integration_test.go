package grpcserver_test

import (
	"context"
	"math/rand"
	"net"
	"path/filepath"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/grpcserver"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
	pb "github.com/zhangchuqi1998/velosearch/proto"
)

const bufSize = 1024 * 1024

func setupServer(t *testing.T) pb.VectorSearchClient {
	t.Helper()
	walPath := filepath.Join(t.TempDir(), "wal.log")
	wal, err := storage.Open(walPath)
	if err != nil {
		t.Fatalf("storage.Open: %v", err)
	}

	lis := bufconn.Listen(bufSize)
	grpcSrv := grpc.NewServer()
	mgr := collection.NewManager()
	pb.RegisterVectorSearchServer(grpcSrv, grpcserver.New(mgr, wal))

	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			t.Logf("Serve exited: %v", err)
		}
	}()

	dial := func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dial),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
		grpcSrv.GracefulStop()
		_ = wal.Close()
	})

	return pb.NewVectorSearchClient(conn)
}

func TestIntegration_EndToEnd(t *testing.T) {
	client := setupServer(t)
	ctx := context.Background()
	const dim = 8
	const n = 100

	if _, err := client.CreateCollection(ctx, &pb.CreateCollectionRequest{
		Name:           "test",
		Dim:            dim,
		Metric:         pb.Metric_METRIC_L2,
		M:              16,
		EfConstruction: 200,
	}); err != nil {
		t.Fatalf("CreateCollection: %v", err)
	}

	rng := rand.New(rand.NewSource(42))
	vecs := make([][]float32, n)
	items := make([]*pb.Item, n)
	for i := 0; i < n; i++ {
		v := make([]float32, dim)
		for j := range v {
			v[j] = rng.Float32()
		}
		vecs[i] = v
		items[i] = &pb.Item{Id: uint32(i), Vector: &pb.Vector{Values: v}}
	}
	insResp, err := client.Insert(ctx, &pb.InsertRequest{Collection: "test", Items: items})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if insResp.Inserted != n {
		t.Fatalf("Inserted: got %d, want %d", insResp.Inserted, n)
	}

	query := vecs[42]
	sResp, err := client.Search(ctx, &pb.SearchRequest{
		Collection: "test",
		Query:      &pb.Vector{Values: query},
		K:          5,
		EfSearch:   50,
	})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(sResp.Hits) == 0 {
		t.Fatalf("Search: empty hits")
	}
	if sResp.Hits[0].Id != 42 {
		t.Errorf("Top hit: got id=%d, want 42 (identical vector should be closest)", sResp.Hits[0].Id)
	}

	delResp, err := client.Delete(ctx, &pb.DeleteRequest{Collection: "test", Ids: []uint32{42}})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if delResp.Deleted != 1 {
		t.Errorf("Deleted: got %d, want 1", delResp.Deleted)
	}

	sResp2, err := client.Search(ctx, &pb.SearchRequest{
		Collection: "test",
		Query:      &pb.Vector{Values: query},
		K:          5,
		EfSearch:   50,
	})
	if err != nil {
		t.Fatalf("Search after delete: %v", err)
	}
	if len(sResp2.Hits) == 0 {
		t.Fatalf("Search after delete: empty hits")
	}
	if sResp2.Hits[0].Id == 42 {
		t.Errorf("Top hit after delete: still id=42, tombstone filter failed")
	}

	stResp, err := client.Stats(ctx, &pb.StatsRequest{Collection: "test"})
	if err != nil {
		t.Fatalf("Stats: %v", err)
	}
	if stResp.NumVectors != n {
		t.Errorf("NumVectors: got %d, want %d", stResp.NumVectors, n)
	}
	if stResp.NumDeleted != 1 {
		t.Errorf("NumDeleted: got %d, want 1", stResp.NumDeleted)
	}
}
