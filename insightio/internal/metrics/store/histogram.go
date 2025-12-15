package store

import "time"

// LatencyHist represents a latency histogram for tracking request latencies
type LatencyHist struct {
	buckets []int64
	counts  []int64
	total   int64
	sumMs   int64
	minMs   int64 // track minimum
	maxMs   int64 // track maximum
}

// NewLatencyHist creates a new latency histogram with the given buckets
func NewLatencyHist(buckets []int64) *LatencyHist {
	return &LatencyHist{
		buckets: buckets,
		counts:  make([]int64, len(buckets)),
		minMs:   -1, // -1 indicates uninitialized
		maxMs:   0,
	}
}

// Observe records a latency observation
func (h *LatencyHist) Observe(d time.Duration) {
	ms := int64(d / time.Millisecond)
	h.total++
	h.sumMs += ms

	// track min/max
	if h.minMs == -1 || ms < h.minMs {
		h.minMs = ms
	}
	if ms > h.maxMs {
		h.maxMs = ms
	}

	for i, upper := range h.buckets {
		if ms <= upper {
			h.counts[i]++
			return
		}
	}

	// last bucket (overflow)
	h.counts[len(h.counts)-1]++
}

// Avg returns the average latency in milliseconds
func (h *LatencyHist) Avg() float64 {
	if h.total == 0 {
		return 0
	}
	return float64(h.sumMs) / float64(h.total)
}

// GetPercentile calculates the percentile latency in milliseconds
func (h *LatencyHist) GetPercentile(percentile float64) float64 {
	if h.total == 0 {
		return 0
	}

	if percentile <= 0 {
		return float64(h.minMs)
	}
	if percentile >= 100 {
		return float64(h.maxMs)
	}

	// Calculate target count for this percentile
	targetCount := float64(h.total) * (percentile / 100.0)

	// Find the bucket that contains this percentile
	accumulated := int64(0)
	for i, count := range h.counts {
		accumulated += count
		if float64(accumulated) >= targetCount {
			// Return the upper bound of this bucket
			return float64(h.buckets[i])
		}
	}

	// Fallback to max
	return float64(h.maxMs)
}

// GetDistribution returns the latency distribution as a map of bucket upper bound to count
func (h *LatencyHist) GetDistribution() map[int64]int64 {
	dist := make(map[int64]int64, len(h.buckets))
	for i, count := range h.counts {
		dist[h.buckets[i]] = count
	}
	return dist
}

// GetMin returns minimum latency in milliseconds
func (h *LatencyHist) GetMin() float64 {
	if h.minMs == -1 {
		return 0
	}
	return float64(h.minMs)
}

// GetMax returns maximum latency in milliseconds
func (h *LatencyHist) GetMax() float64 {
	return float64(h.maxMs)
}

// GetMedian returns median (p50) latency in milliseconds
func (h *LatencyHist) GetMedian() float64 {
	return h.GetPercentile(50)
}
