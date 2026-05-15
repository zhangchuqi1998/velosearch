// VeloSearch gRPC server entrypoint.
package main

import (
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/grpcserver"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
	pb "github.com/zhangchuqi1998/velosearch/proto"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC listen address")
	metricsAddr := flag.String("metrics-addr", ":9090", "HTTP listen address for Prometheus /metrics (empty disables)")
	dataDir := flag.String("data-dir", "./data", "WAL data directory")
	walEnabled := flag.Bool("wal", true, "enable write-ahead log (disable for benchmarks)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mgr := collection.NewManager()
	var wal *storage.WAL

	if *walEnabled {
		if err := os.MkdirAll(*dataDir, 0755); err != nil {
			slog.Error("mkdir failed", "err", err, "data_dir", *dataDir)
			os.Exit(1)
		}
		walPath := filepath.Join(*dataDir, "wal.log")

		slog.Info("replaying WAL", "path", walPath)
		nRec := 0
		if err := storage.Replay(walPath, func(rec *storage.Record) error {
			nRec++
			return ApplyWALRecord(mgr, rec)
		}); err != nil {
			slog.Error("replay failed", "err", err)
			os.Exit(1)
		}
		slog.Info("replay done", "records", nRec)

		w, err := storage.Open(walPath)
		if err != nil {
			slog.Error("open wal failed", "err", err)
			os.Exit(1)
		}
		wal = w
		defer wal.Close()
	} else {
		slog.Warn("WAL disabled — durability off, restarts will lose all data")
	}

	srv := grpcserver.New(mgr, wal)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		slog.Error("listen failed", "err", err)
		os.Exit(1)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterVectorSearchServer(grpcSrv, srv)

	// Prometheus /metrics on a separate port so scrapers don't compete with
	// gRPC traffic. Use a dedicated ServeMux so we don't pollute the global
	// http.DefaultServeMux (which third-party libs sometimes register on).
	if *metricsAddr != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		metricsSrv := &http.Server{Addr: *metricsAddr, Handler: mux}
		go func() {
			slog.Info("metrics listening", "addr", *metricsAddr)
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("metrics serve failed", "err", err)
			}
		}()
	}

	go func() {
		slog.Info("server listening", "addr", *addr)
		if err := grpcSrv.Serve(lis); err != nil {
			slog.Error("serve failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	slog.Info("shutting down")
	grpcSrv.GracefulStop()
}
