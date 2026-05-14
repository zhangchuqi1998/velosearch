// crash_client is a deterministic-workload driver for the Day 11 crash
// recovery test.
//
//	go run ./benchmark/crash_client -addr=localhost:50052 -mode=write  -n=1000
//	go run ./benchmark/crash_client -addr=localhost:50052 -mode=verify -n=1000
//
// In write mode it creates the "crash" collection and inserts n vectors
// where vec[i][j] = float32(i+j). The function is pure so verify mode can
// rebuild the same vector for each id without any shared state.
//
// In verify mode it Search(k=1) for each id with that same vector; the
// nearest hit must be id itself (identical-vector distance is zero), or
// the recovery is considered failed and the process exits with code 1.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/zhangchuqi1998/velosearch/proto"
)

const (
	collName  = "crash"
	dim       = 32
	batchSize = 100
)

func buildVec(i int) []float32 {
	v := make([]float32, dim)
	for j := 0; j < dim; j++ {
		v[j] = float32(i + j)
	}
	return v
}

func main() {
	addr := flag.String("addr", "localhost:50052", "gRPC server address")
	mode := flag.String("mode", "", "write or verify")
	n := flag.Int("n", 1000, "number of vectors")
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewVectorSearchClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	switch *mode {
	case "write":
		writeMode(ctx, client, *n)
	case "verify":
		verifyMode(ctx, client, *n)
	default:
		log.Fatalf("mode must be 'write' or 'verify', got %q", *mode)
	}
}

func writeMode(ctx context.Context, client pb.VectorSearchClient, n int) {
	_, err := client.CreateCollection(ctx, &pb.CreateCollectionRequest{
		Name:           collName,
		Dim:            int32(dim),
		Metric:         pb.Metric_METRIC_L2,
		M:              16,
		EfConstruction: 200,
	})
	if err != nil && !isAlreadyExists(err) {
		log.Fatalf("CreateCollection: %v", err)
	}

	for start := 0; start < n; start += batchSize {
		end := start + batchSize
		if end > n {
			end = n
		}
		items := make([]*pb.Item, end-start)
		for i := start; i < end; i++ {
			items[i-start] = &pb.Item{
				Id:     uint32(i),
				Vector: &pb.Vector{Values: buildVec(i)},
			}
		}
		if _, err := client.Insert(ctx, &pb.InsertRequest{Collection: collName, Items: items}); err != nil {
			log.Fatalf("Insert [%d, %d): %v", start, end, err)
		}
	}
	fmt.Printf("DONE writing %d vectors\n", n)
}

func verifyMode(ctx context.Context, client pb.VectorSearchClient, n int) {
	misses := 0
	firstMiss := -1
	for i := 0; i < n; i++ {
		resp, err := client.Search(ctx, &pb.SearchRequest{
			Collection: collName,
			Query:      &pb.Vector{Values: buildVec(i)},
			K:          1,
			EfSearch:   50,
		})
		if err != nil {
			log.Fatalf("Search id=%d: %v", i, err)
		}
		if len(resp.Hits) == 0 || resp.Hits[0].Id != uint32(i) {
			misses++
			if firstMiss == -1 {
				firstMiss = i
			}
		}
	}
	fmt.Printf("verified n=%d misses=%d", n, misses)
	if misses > 0 {
		fmt.Printf(" first_miss=%d\n", firstMiss)
		os.Exit(1)
	}
	fmt.Println()
}

func isAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	// crude fallback: gRPC code is "AlreadyExists", which serializes into the error message
	if errors.Is(err, context.Canceled) {
		return false
	}
	return strings.Contains(err.Error(), "AlreadyExists") || strings.Contains(err.Error(), "already exists")
}
