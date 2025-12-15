package auth

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	apiKeyHeader = "x-api-key"
)

var (
	ErrMissingAPIKey = errors.New("API key is missing")
	ErrInvalidAPIKey = errors.New("API key is invalid")
)

// APIKeyValidator validates API keys for gRPC requests
type APIKeyValidator struct {
	mu        sync.RWMutex
	validKeys map[string]bool
}

// NewAPIKeyValidator creates a new API key validator with the given valid keys
func NewAPIKeyValidator(validKeys []string) *APIKeyValidator {
	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		keyMap[key] = true
	}
	return &APIKeyValidator{
		validKeys: keyMap,
	}
}

// AddKey adds a new valid API key
func (v *APIKeyValidator) AddKey(key string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.validKeys[key] = true
}

// RemoveKey removes an API key
func (v *APIKeyValidator) RemoveKey(key string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.validKeys, key)
}

// ValidateAPIKey validates the API key from the gRPC context metadata
func (v *APIKeyValidator) ValidateAPIKey(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, ErrMissingAPIKey.Error())
	}

	apiKeys := md.Get(apiKeyHeader)
	if len(apiKeys) == 0 || apiKeys[0] == "" {
		return status.Error(codes.Unauthenticated, ErrMissingAPIKey.Error())
	}

	apiKey := apiKeys[0]

	v.mu.RLock()
	defer v.mu.RUnlock()

	if !v.validKeys[apiKey] {
		return status.Error(codes.Unauthenticated, ErrInvalidAPIKey.Error())
	}

	return nil
}
