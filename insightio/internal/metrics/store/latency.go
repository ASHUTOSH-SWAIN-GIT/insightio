package store

import "time"

// RecordLatency records latency for a specific method
func (m *MetricStore) RecordLatency(method string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	hist, ok := m.latencyMap[method]
	if !ok {
		hist = NewLatencyHist(m.buckets)
		m.latencyMap[method] = hist
	}

	hist.Observe(d)
}

// GetLatencyPercentile returns the latency percentile for a specific method
func (m *MetricStore) GetLatencyPercentile(method string, percentile float64) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hist, ok := m.latencyMap[method]
	if !ok {
		return 0
	}

	return hist.GetPercentile(percentile)
}

// GetLatencyDistribution returns the latency distribution for a specific method
func (m *MetricStore) GetLatencyDistribution(method string) map[int64]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hist, ok := m.latencyMap[method]
	if !ok {
		return make(map[int64]int64)
	}

	return hist.GetDistribution()
}

// LatencyStats represents comprehensive latency statistics
type LatencyStats struct {
	Min       float64
	Max       float64
	Avg       float64
	Median    float64
	P95       float64
	P99       float64
	TotalReqs int64
}

// GetLatencyStats returns comprehensive latency stats for a method
func (m *MetricStore) GetLatencyStats(method string) LatencyStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hist, ok := m.latencyMap[method]
	if !ok {
		return LatencyStats{}
	}

	return LatencyStats{
		Min:       hist.GetMin(),
		Max:       hist.GetMax(),
		Avg:       hist.Avg(),
		Median:    hist.GetMedian(),
		P95:       hist.GetPercentile(95),
		P99:       hist.GetPercentile(99),
		TotalReqs: hist.total,
	}
}
