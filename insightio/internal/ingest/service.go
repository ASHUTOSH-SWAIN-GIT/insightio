package ingest

import (
	"context"
	"io"
	"log"

	pb "github.com/ASHUTOSH-SWAIN-GIT/insightio/proto"
)

// IngestServiceServer implements the gRPC ingest service.
// It pushes incoming events into eventChan for the worker to process.
type IngestServiceServer struct {
	pb.UnimplementedIngestServiceServer
	eventChan chan<- *pb.Event
}

// NewIngestService creates a new ingest service with the event channel.
func NewIngestService(eventChan chan<- *pb.Event) *IngestServiceServer {
	return &IngestServiceServer{
		eventChan: eventChan,
	}
}

// SendEvent handles unary RPC for a single event.
func (s *IngestServiceServer) SendEvent(ctx context.Context, event *pb.Event) (*pb.Ack, error) {

	if event == nil || event.Type == "" {
		return &pb.Ack{
			Ok:      false,
			Message: "Event or event type is missing",
		}, nil
	}

	// Push event to worker channel
	s.eventChan <- event

	return &pb.Ack{
		Ok:      true,
		Message: "Event received successfully",
	}, nil
}

// SendEventStream handles client-streaming RPC where multiple events are sent.
func (s *IngestServiceServer) SendEventStream(stream pb.IngestService_SendEventStreamServer) error {

	count := 0

	for {
		event, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received %d events in stream", count)

			// Send final ack and close stream
			return stream.SendAndClose(&pb.Ack{
				Ok:      true,
				Message: "Stream received successfully",
			})
		}

		if err != nil {
			log.Printf("Stream error: %v", err)
			return err
		}

		// Validate event
		if event.Type == "" {
			log.Println("Skipping event with missing type")
			continue
		}

		// Push event to worker channel
		s.eventChan <- event
		count++
	}
}
