package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/auth"
	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/config"
	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/ingest"
	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics/store"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics"
	pb "github.com/ASHUTOSH-SWAIN-GIT/insightio/proto"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Create metric store with configured window size
	store := store.NewMetricStore(cfg.MetricsWindow)

	eventChan := make(chan *pb.Event, 1000) // channel use to send events from the grpc service to worker

	// Start worker that processes events and updates the store
	worker := ingest.NewWorker(eventChan, store)
	worker.Start()

	// Initialize API key validator with keys from config
	validator := auth.NewAPIKeyValidator(cfg.APIKeys)

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor(store, validator)),
		grpc.StreamInterceptor(metrics.StreamServerInterceptor(store, validator)),
	)

	// Register services
	pb.RegisterIngestServiceServer(
		grpcServer,
		ingest.NewIngestService(eventChan),
	)
	pb.RegisterMetricsServiceServer(
		grpcServer,
		metrics.NewMetricsService(store),
	)

	// Start TCP listener on configured port
	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", cfg.GRPCPort, err)
	}

	log.Printf("InsightIO analytics engine running on port %d", cfg.GRPCPort)
	log.Printf("Metrics window: %d seconds", cfg.MetricsWindow)
	log.Printf("Environment: %s", cfg.Env)
	log.Printf("API key validation enabled (%d key(s) configured)", len(cfg.APIKeys))

	// Start the server loop
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
