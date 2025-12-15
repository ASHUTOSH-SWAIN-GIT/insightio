package metrics

import (
	"time"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/auth"
	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics/store"

	"google.golang.org/grpc"
)

// StreamServerInterceptor wraps streaming RPC calls and validates API keys
func StreamServerInterceptor(store *store.MetricStore, validator *auth.APIKeyValidator) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		start := time.Now()

		// Validate API key before processing request
		if validator != nil {
			if err := validator.ValidateAPIKey(ss.Context()); err != nil {
				// Record failed auth attempt
				store.RecordRequest(info.FullMethod)
				store.RecordError(info.FullMethod)
				return err
			}
		}

		err := handler(srv, ss)

		elapsed := time.Since(start)
		method := info.FullMethod

		store.RecordRequest(method)
		store.RecordLatency(method, elapsed)
		if err != nil {
			store.RecordError(method)
		}

		return err
	}
}
