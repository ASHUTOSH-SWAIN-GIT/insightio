package metrics

import (
	"context"
	"time"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/auth"
	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics/store"

	"google.golang.org/grpc"
)

// wraps the rpc calls to track latency, errors, throughput, and validate API keys
func UnaryServerInterceptor(store *store.MetricStore, validator *auth.APIKeyValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Validate API key before processing request
		if validator != nil {
			if err := validator.ValidateAPIKey(ctx); err != nil {
				// Record failed auth attempt
				store.RecordRequest(info.FullMethod)
				store.RecordError(info.FullMethod)
				return nil, err
			}
		}

		// Process request
		resp, err := handler(ctx, req)

		elapsed := time.Since(start)
		method := info.FullMethod

		store.RecordRequest(method)
		store.RecordLatency(method, elapsed)
		if err != nil {
			store.RecordError(method)
		}

		return resp, err
	}
}
