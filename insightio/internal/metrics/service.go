package metrics

import (
	"context"
	"log"
	"time"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics/store"
	pb "github.com/ASHUTOSH-SWAIN-GIT/insightio/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MetricsServiceServer implements the metrics service.
type MetricsServiceServer struct {
	pb.UnimplementedMetricsServiceServer
	store *store.MetricStore
}

// NewMetricsService returns a new metrics service instance.
func NewMetricsService(store *store.MetricStore) *MetricsServiceServer {
	return &MetricsServiceServer{store: store}
}

// GetMetrics returns a snapshot of requested metrics.
func (s *MetricsServiceServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.MetricResponse, error) {

	resp := &pb.MetricResponse{}

	// if client requested no specific metrics, return defaults
	if len(req.MetricsNames) == 0 {
		resp.Metrics = append(resp.Metrics,
			s.makeMetric("total_events", float64(s.store.GetTotalEvents())),
			s.makeMetric("events_per_window", float64(s.store.GetEventsPerWindow())),
			s.makeMetric("total_throughput", s.store.GetTotalThroughput()),
			s.makeMetric("total_error_rate", s.store.GetTotalErrorRate()),
		)
		return resp, nil
	}

	// return only requested metrics
	for _, name := range req.MetricsNames {

		switch name {

		case "total_events":
			resp.Metrics = append(resp.Metrics,
				s.makeMetric("total_events", float64(s.store.GetTotalEvents())),
			)

		case "events_per_window":
			resp.Metrics = append(resp.Metrics,
				s.makeMetric("events_per_window", float64(s.store.GetEventsPerWindow())),
			)

		case "total_throughput":
			resp.Metrics = append(resp.Metrics,
				s.makeMetric("total_throughput", s.store.GetTotalThroughput()),
			)

		case "total_error_rate":
			resp.Metrics = append(resp.Metrics,
				s.makeMetric("total_error_rate", s.store.GetTotalErrorRate()),
			)

		default:
			log.Printf("Unknown metric requested: %s", name)
		}
	}

	return resp, nil
}

// SubscribeMetrics streams metrics in real time.
func (s *MetricsServiceServer) SubscribeMetrics(req *pb.GetMetricsRequest, stream pb.MetricsService_SubscribeMetricsServer) error {

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			// Send multiple metrics in the stream
			metrics := []*pb.Metric{
				s.makeMetric("events_per_window", float64(s.store.GetEventsPerWindow())),
				s.makeMetric("total_throughput", s.store.GetTotalThroughput()),
				s.makeMetric("total_error_rate", s.store.GetTotalErrorRate()),
			}

			for _, metric := range metrics {
				if err := stream.Send(metric); err != nil {
					log.Println("Client disconnected:", err)
					return err
				}
			}

		case <-stream.Context().Done():
			log.Println("Client unsubscribed from stream")
			return nil
		}
	}
}

// Helper: create metric with timestamp
func (s *MetricsServiceServer) makeMetric(name string, value float64) *pb.Metric {
	return &pb.Metric{
		Name:      name,
		Value:     value,
		Timestamp: timestamppb.Now(),
	}
}
