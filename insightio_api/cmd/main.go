package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ASHUTOSH-SWAIN-GIT/insightio/proto"
)

const (
	httpPort     = ":8080"
	grpcTarget   = "localhost:50051"
	ingestPath   = "/v1/event"
	apiKeyHeader = "x-api-key"
)

// mirrors the structure recieved from the client via http/json
type IngestPayload struct {
	Type     string            `json:"type"`
	UserId   string            `json:"user_id,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Value    float64           `json:"value,omitempty"`
}

var gRPCClient pb.IngestServiceClient

func main() {
	//persistent grpc client  connection to the insightio server
	conn, err := grpc.Dial(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server %s: %v", grpcTarget, err)
	}
	defer conn.Close()

	//grpc service client  interface
	gRPCClient = pb.NewIngestServiceClient(conn)

	//http router
	http.HandleFunc(ingestPath, ingestHandler)

	log.Printf("Starting API Gateway on %s, proxying to %s", httpPort, grpcTarget)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal(err)
	}
}

// ingestHandler processes the incoming HTTP/JSON request and proxies it to gRPC.
func ingestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size to prevent DoS (10MB max)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var payload IngestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if payload.Type == "" {
		http.Error(w, "Field 'type' is required", http.StatusBadRequest)
		return
	}

	// 1. Transcode and Enrich to Protobuf Message
	event := &pb.Event{
		Id:        uuid.New().String(),
		Type:      payload.Type,
		Value:     payload.Value,
		UserId:    payload.UserId,
		Metadata:  payload.Metadata,
		Timestamp: timestamppb.Now(),
	}

	// 2. Extract API Key and propagate it via gRPC Metadata
	apiKey := r.Header.Get(apiKeyHeader)
	if apiKey == "" {
		http.Error(w, apiKeyHeader+" header missing", http.StatusUnauthorized)
		return
	}

	// Use request context instead of Background() for proper cancellation
	ctx := r.Context()
	// Append the API Key to the gRPC context for the server's Interceptor
	ctx = metadata.AppendToOutgoingContext(ctx, apiKeyHeader, apiKey)

	// 3. Call the gRPC Backend (Unary RPC)
	gRPCContext, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	resp, err := gRPCClient.SendEvent(gRPCContext, event)

	// 4. Handle response from gRPC
	if err != nil {
		// gRPC error (e.g., connection issue, deadline exceeded)
		log.Printf("gRPC SendEvent failed: %v", err)
		// Don't expose internal error details to clients
		http.Error(w, "Backend service unavailable", http.StatusServiceUnavailable)
		return
	}

	if !resp.Ok {
		// Logical rejection from IngestService (e.g., event type missing)
		http.Error(w, "Ingestion rejected: "+resp.Message, http.StatusBadRequest)
		return
	}

	// 5. Success: Event accepted (queued to worker channel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // HTTP 202: Accepted for asynchronous processing
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": resp.Message,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
