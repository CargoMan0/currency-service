package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Counter *prometheus.CounterVec
	Latency *prometheus.HistogramVec
}

func Init() (*Metrics, error) {
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "currency_requests_total",
			Help: "Total number of requests handled by the currency service",
		},
		[]string{"method"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "currency_request_duration_seconds",
			Help:    "Histogram of response times for requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	err := prometheus.Register(requestCount)
	if err != nil {
		return nil, fmt.Errorf("error registering requests metrics: %w", err)
	}
	err = prometheus.Register(requestDuration)
	if err != nil {
		return nil, fmt.Errorf("error registering requests duration metrics: %w", err)
	}

	return &Metrics{
		Counter: requestCount,
		Latency: requestDuration,
	}, nil
}
