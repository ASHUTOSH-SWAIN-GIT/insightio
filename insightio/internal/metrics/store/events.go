package store

import "time"

// AddEvent records an event in the store
func (m *MetricStore) AddEvent(eventType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalEvents++
	now := time.Now()
	m.eventTypeCounts[eventType]++
	m.eventTimestamps = append(m.eventTimestamps, time.Now())

	//clean up the timestamps  outside the window
	threshold := now.Add(-m.windowSize)
	validTimestamps := make([]time.Time, 0, len(m.eventTimestamps))
	for _, ts := range m.eventTimestamps {
		if ts.After(threshold) {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	m.eventTimestamps = validTimestamps
}

// GetTotalEvents returns the total number of events recorded
func (m *MetricStore) GetTotalEvents() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalEvents
}

// GetEventTypeCount returns the count of events for a specific type
func (m *MetricStore) GetEventTypeCount(eventType string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.eventTypeCounts[eventType]
}

// GetEventsPerWindow returns the number of events within the sliding window
func (m *MetricStore) GetEventsPerWindow() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	threshold := now.Add(-m.windowSize)

	count := int64(0)
	for _, ts := range m.eventTimestamps {
		if ts.After(threshold) {
			count++
		}
	}
	return count
}
