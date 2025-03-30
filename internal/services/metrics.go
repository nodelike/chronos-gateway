package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsCollector struct {
	RequestCounter     *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
	LocationCounter    *prometheus.CounterVec
	LocationLatency    *prometheus.HistogramVec
	BatchSizeHistogram *prometheus.HistogramVec
}

func NewMetricsCollector() *MetricsCollector {
	// Use promauto which automatically registers the metrics
	collector := &MetricsCollector{
		RequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "Total API requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_request_duration_seconds",
				Help:    "API request duration distribution",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		LocationCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "location_events_total",
				Help: "Total location events received",
			},
			[]string{"event_type", "source"},
		),
		LocationLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "location_event_latency_seconds",
				Help:    "Time difference between event creation and reception",
				Buckets: []float64{0.1, 0.5, 1, 5, 10, 30, 60, 120, 300, 600},
			},
			[]string{"event_type"},
		),
		BatchSizeHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "location_batch_size",
				Help:    "Size of location event batches",
				Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"source"},
		),
	}

	// No need to register metrics manually since promauto does it for us
	return collector
}

// RecordLocationEvent records metrics for a single location event
func (m *MetricsCollector) RecordLocationEvent(eventType, source string, latencySeconds float64) {
	m.LocationCounter.WithLabelValues(eventType, source).Inc()
	m.LocationLatency.WithLabelValues(eventType).Observe(latencySeconds)
}

// RecordBatchSize records the size of a batch of location events
func (m *MetricsCollector) RecordBatchSize(source string, size int) {
	m.BatchSizeHistogram.WithLabelValues(source).Observe(float64(size))
}
