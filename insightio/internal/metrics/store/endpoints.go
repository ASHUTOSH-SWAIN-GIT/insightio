package store

import "sort"

// EndpointStats represents statistics for an endpoint
type EndpointStats struct {
	Method string
	AvgMs  float64
	Reqs   int64
	Errs   int64
}

// GetTopSlowestEndpoints returns the k slowest endpoints sorted by average latency
func (m *MetricStore) GetTopSlowestEndpoints(k int) []EndpointStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type entry struct {
		method string
		avgMs  float64
		reqs   int64
		errs   int64
	}

	list := []entry{}
	for method, hist := range m.latencyMap {
		list = append(list, entry{
			method: method,
			avgMs:  hist.Avg(),
			reqs:   m.reqCount[method],
			errs:   m.errCount[method],
		})
	}

	// sort by avg latency descending
	sort.Slice(list, func(i, j int) bool {
		return list[i].avgMs > list[j].avgMs
	})

	if k > len(list) {
		k = len(list)
	}

	out := make([]EndpointStats, k)
	for i := 0; i < k; i++ {
		out[i] = EndpointStats{
			Method: list[i].method,
			AvgMs:  list[i].avgMs,
			Reqs:   list[i].reqs,
			Errs:   list[i].errs,
		}
	}

	return out
}
