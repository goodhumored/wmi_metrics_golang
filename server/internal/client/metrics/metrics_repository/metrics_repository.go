package metrics_repository

import (
	"goodhumored/wmi-metrics-server/internal/client/metrics"
)

type MemoryMetricsRepository struct {
	metricsStorage map[string][]metrics.Metrics
}

func NewMemoryMetricsRepository() *MemoryMetricsRepository {
	return &MemoryMetricsRepository{
		metricsStorage: make(map[string][]metrics.Metrics),
	}
}

func (r *MemoryMetricsRepository) StoreMetrics(clientID string, m metrics.Metrics) error {
	r.metricsStorage[clientID] = append(r.metricsStorage[clientID], m)
	return nil
}

func (r *MemoryMetricsRepository) GetLatestMetrics(clientID string) (metrics.Metrics, bool) {
	metricsList, exists := r.metricsStorage[clientID]
	if !exists || len(metricsList) == 0 {
		return metrics.Metrics{}, false
	}
	return metricsList[len(metricsList)-1], true
}

func (r *MemoryMetricsRepository) GetMetricsForPeriod(clientID string, from, to int64) ([]metrics.Metrics, error) {
	metricsList, exists := r.metricsStorage[clientID]
	if !exists {
		return nil, nil
	}
	left, right := 0, len(metricsList)-1
	for left <= right {
		mid := left + (right-left)/2
		if metricsList[mid].Timestamp < from {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	startIndex := left
	left, right = 0, len(metricsList)-1
	for left <= right {
		mid := left + (right-left)/2
		if metricsList[mid].Timestamp <= to {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	endIndex := right
	if startIndex > endIndex {
		return []metrics.Metrics{}, nil
	}
	if endIndex >= len(metricsList) {
		endIndex = len(metricsList) - 1
	}
	return metricsList[startIndex : endIndex+1], nil
}
