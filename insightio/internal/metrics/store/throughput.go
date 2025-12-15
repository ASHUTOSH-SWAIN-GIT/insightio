package store

import "time"

// RecordRequest records a request for a specific method
func (m *MetricStore) RecordRequest(method string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reqCount[method]++
	// Track timestamp for throughput calculation
	now := time.Now()
	m.reqTimestamps[method] = append(m.reqTimestamps[method], now)

	// Clean up old timestamps outside the window to prevent memory growth
	threshold := now.Add(-m.windowSize)
	validTimestamps := make([]time.Time, 0, len(m.reqTimestamps[method]))
	for _, ts := range m.reqTimestamps[method] {
		if ts.After(threshold) {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	m.reqTimestamps[method] = validTimestamps
}

// GetThroughput returns requests per second for a specific method within the window
func (m *MetricStore) GetThroughput(method string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	timestamps, ok := m.reqTimestamps[method]
	if !ok || len(timestamps) == 0 {
		return 0
	}

	// Count requests in the window
	now := time.Now()
	threshold := now.Add(-m.windowSize)
	count := int64(0)
	for _, ts := range timestamps {
		if ts.After(threshold) {
			count++
		}
	}

	// Calculate requests per second
	windowSeconds := m.windowSize.Seconds()
	if windowSeconds <= 0 {
		return 0
	}
	return float64(count) / windowSeconds
}

// GetTotalThroughput returns overall requests per second across all methods
func (m *MetricStore) GetTotalThroughput() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	threshold := now.Add(-m.windowSize)
	totalCount := int64(0)

	for _, timestamps := range m.reqTimestamps {
		for _, ts := range timestamps {
			if ts.After(threshold) {
				totalCount++
			}
		}
	}

	windowSeconds := m.windowSize.Seconds()
	if windowSeconds <= 0 {
		return 0
	}
	return float64(totalCount) / windowSeconds
}
