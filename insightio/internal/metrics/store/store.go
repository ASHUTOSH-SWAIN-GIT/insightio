package store

import (
	"sync"
	"time"
)

// Histogram buckets in milliseconds
var DefaultBuckets = []int64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000}

// MetricStore is the main metrics storage structure
type MetricStore struct {
	mu              sync.RWMutex
	totalEvents     int64
	eventTypeCounts map[string]int64
	eventTimestamps []time.Time
	windowSize      time.Duration
	reqCount        map[string]int64
	errCount        map[string]int64
	latencyMap      map[string]*LatencyHist
	buckets         []int64
	reqTimestamps   map[string][]time.Time // method -> timestamps
}

// NewMetricStore creates a new metrics store with the specified window size in seconds
func NewMetricStore(windowSeconds int) *MetricStore {
	return &MetricStore{
		eventTypeCounts: make(map[string]int64),
		eventTimestamps: make([]time.Time, 0, 1024),
		windowSize:      time.Duration(windowSeconds) * time.Second,

		reqCount:      make(map[string]int64),
		errCount:      make(map[string]int64),
		latencyMap:    make(map[string]*LatencyHist),
		buckets:       DefaultBuckets,
		reqTimestamps: make(map[string][]time.Time),
	}
}
