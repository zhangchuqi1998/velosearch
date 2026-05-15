package grpcserver

import (
	"context"
	"errors"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
	pb "github.com/zhangchuqi1998/velosearch/proto"
)

type Server struct {
	pb.UnimplementedVectorSearchServer
	mgr *collection.Manager
	wal *storage.WAL
}

func New(mgr *collection.Manager, wal *storage.WAL) *Server {
	return &Server{mgr: mgr, wal: wal}
}

func metricFromProto(m pb.Metric) (collection.Metric, error) {
	switch m {
	case pb.Metric_METRIC_L2:
		return collection.MetricL2, nil
	case pb.Metric_METRIC_COSINE:
		return collection.MetricCosine, nil
	default:
		return 0, status.Error(codes.InvalidArgument, "metric must be METRIC_L2 or METRIC_COSINE")
	}
}

func mapMgrErr(err error) error {
	switch {
	case errors.Is(err, collection.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, collection.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, collection.ErrDimMismatch):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *Server) CreateCollection(ctx context.Context, req *pb.CreateCollectionRequest) (*pb.CreateCollectionResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Dim <= 0 {
		return nil, status.Error(codes.InvalidArgument, "dim must be > 0")
	}
	if req.M < 4 {
		return nil, status.Error(codes.InvalidArgument, "m must be >= 4")
	}
	if req.EfConstruction < req.M {
		return nil, status.Error(codes.InvalidArgument, "ef_construction must be >= m")
	}
	metric, err := metricFromProto(req.Metric)
	if err != nil {
		return nil, err
	}

	rec := &storage.Record{Op: storage.OpCreateColl, CreateColl: &storage.CreateColl{
		Name: req.Name, Dim: req.Dim, Metric: uint8(metric),
		M: req.M, EfConstruction: req.EfConstruction,
	}}
	if err := s.wal.Append(rec); err != nil {
		return nil, status.Error(codes.Internal, "wal append: "+err.Error())
	}

	cfg := collection.Config{
		Name:           req.Name,
		Dim:            int(req.Dim),
		Metric:         metric,
		M:              int(req.M),
		EfConstruction: int(req.EfConstruction),
	}
	if err := s.mgr.Create(cfg); err != nil {
		return nil, mapMgrErr(err)
	}
	slog.Info("CreateCollection", "name", req.Name, "dim", req.Dim, "m", req.M, "ef_construction", req.EfConstruction)
	return &pb.CreateCollectionResponse{Created: true}, nil
}

func (s *Server) DropCollection(ctx context.Context, req *pb.DropCollectionRequest) (*pb.DropCollectionResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	rec := &storage.Record{Op: storage.OpDropColl, DropColl: &storage.DropColl{Name: req.Name}}
	if err := s.wal.Append(rec); err != nil {
		return nil, status.Error(codes.Internal, "wal append: "+err.Error())
	}

	if err := s.mgr.Drop(req.Name); err != nil {
		return nil, mapMgrErr(err)
	}
	slog.Info("DropCollection", "name", req.Name)
	return &pb.DropCollectionResponse{Dropped: true}, nil
}

func (s *Server) ListCollections(ctx context.Context, req *pb.ListCollectionsRequest) (*pb.ListCollectionsResponse, error) {
	names := s.mgr.List()
	slog.Info("ListCollections", "count", len(names))
	return &pb.ListCollectionsResponse{Names: names}, nil
}

func (s *Server) Insert(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection is required")
	}
	col, err := s.mgr.Get(req.Collection)
	if err != nil {
		return nil, mapMgrErr(err)
	}
	dim := col.Config.Dim
	for i, item := range req.Items {
		if item.Vector == nil || len(item.Vector.Values) != dim {
			return nil, status.Errorf(codes.InvalidArgument,
				"items[%d]: vector dim mismatch (got %d, want %d)", i, len(item.GetVector().GetValues()), dim)
		}
	}
	for _, item := range req.Items {
		rec := &storage.Record{Op: storage.OpInsert, Insert: &storage.Insert{
			Collection: req.Collection, ID: item.Id, Vector: item.Vector.Values,
		}}
		if err := s.wal.Append(rec); err != nil {
			return nil, status.Error(codes.Internal, "wal append: "+err.Error())
		}
		col.Index.Insert(item.Id, item.Vector.Values)
	}
	InsertCounter.WithLabelValues(req.Collection).Add(float64(len(req.Items)))
	CollectionSize.WithLabelValues(req.Collection).Set(float64(col.Index.SnapshotStats().NumVectors))
	slog.Info("Insert", "collection", req.Collection, "count", len(req.Items))
	return &pb.InsertResponse{Inserted: int32(len(req.Items))}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection is required")
	}
	col, err := s.mgr.Get(req.Collection)
	if err != nil {
		return nil, mapMgrErr(err)
	}
	deleted := 0
	for _, id := range req.Ids {
		rec := &storage.Record{Op: storage.OpDelete, Delete: &storage.Delete{
			Collection: req.Collection, ID: id,
		}}
		if err := s.wal.Append(rec); err != nil {
			return nil, status.Error(codes.Internal, "wal append: "+err.Error())
		}
		if err := col.Index.Delete(id); err == nil {
			deleted++
		}
	}
	DeleteCounter.WithLabelValues(req.Collection).Add(float64(deleted))
	slog.Info("Delete", "collection", req.Collection, "requested", len(req.Ids), "deleted", deleted)
	return &pb.DeleteResponse{Deleted: int32(deleted)}, nil
}

func (s *Server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	timer := prometheus.NewTimer(SearchLatency.WithLabelValues(req.Collection))
	defer timer.ObserveDuration()

	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection is required")
	}
	if req.K <= 0 {
		return nil, status.Error(codes.InvalidArgument, "k must be > 0")
	}
	if req.EfSearch < req.K {
		return nil, status.Error(codes.InvalidArgument, "ef_search must be >= k")
	}
	if req.Query == nil {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}
	col, err := s.mgr.Get(req.Collection)
	if err != nil {
		return nil, mapMgrErr(err)
	}
	if len(req.Query.Values) != col.Config.Dim {
		return nil, status.Errorf(codes.InvalidArgument,
			"query dim mismatch (got %d, want %d)", len(req.Query.Values), col.Config.Dim)
	}
	cands := col.Index.Search(req.Query.Values, int(req.K), int(req.EfSearch))
	hits := make([]*pb.Hit, len(cands))
	for i, c := range cands {
		hits[i] = &pb.Hit{Id: c.ID, Distance: c.Dist}
	}
	slog.Info("Search", "collection", req.Collection, "k", req.K, "ef_search", req.EfSearch, "hits", len(hits))
	return &pb.SearchResponse{Hits: hits}, nil
}

func (s *Server) Stats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection is required")
	}
	col, err := s.mgr.Get(req.Collection)
	if err != nil {
		return nil, mapMgrErr(err)
	}
	st := col.Index.SnapshotStats()
	slog.Info("Stats", "collection", req.Collection, "num_vectors", st.NumVectors)
	return &pb.StatsResponse{
		NumVectors: int32(st.NumVectors),
		NumDeleted: int32(st.NumDeleted),
		NumLayers:  int32(st.NumLayers),
		MemBytes:   st.MemBytes,
	}, nil
}
