package store

// RecordError records an error for a specific method
func (m *MetricStore) RecordError(method string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errCount[method]++
}

// GetErrorRate returns error rate as a percentage (0-100) for a specific method
func (m *MetricStore) GetErrorRate(method string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reqs := m.reqCount[method]
	if reqs == 0 {
		return 0
	}

	errs := m.errCount[method]
	return (float64(errs) / float64(reqs)) * 100.0
}

// GetTotalErrorRate returns overall error rate as a percentage across all methods
func (m *MetricStore) GetTotalErrorRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalReqs := int64(0)
	totalErrs := int64(0)

	for method := range m.reqCount {
		totalReqs += m.reqCount[method]
		totalErrs += m.errCount[method]
	}

	if totalReqs == 0 {
		return 0
	}

	return (float64(totalErrs) / float64(totalReqs)) * 100.0
}
