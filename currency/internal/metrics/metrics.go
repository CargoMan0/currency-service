package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Counter   *prometheus.CounterVec
	Latency   *prometheus.HistogramVec
	AppUptime prometheus.Gauge
}

func Init() *Metrics {
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

	appUptime := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "currency_service_uptime_seconds",
			Help: "Time since service start in seconds",
		},
	)
	return &Metrics{
		Counter:   requestCount,
		Latency:   requestDuration,
		AppUptime: appUptime,
	}
}
