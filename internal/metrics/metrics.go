package metrics

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	http.Handler
	registry *prometheus.Registry
	JobPosts prometheus.Gauge

	RequestsTotal    prometheus.CounterVec
	RequestsDuration prometheus.HistogramVec
}

// NewMetricsSubsystem returns the metrics endpoint HTTP handler and the Prometheus metrics collectors for the server and middleware.
func NewMetricsSubsystem(db *db.Queries) *PrometheusMetrics {
	metricsBackend := newPrometheusMetrics()
	metricsBackend.Handler = promhttp.HandlerFor(metricsBackend.registry, promhttp.HandlerOpts{})
	ticker := time.NewTicker(120 * time.Second)
	go func() {
		for ; ; <-ticker.C {
			jobPosts, err := db.ListJobPosts(context.Background())
			if err != nil {
				log.Println("error listing job posts:", err)
				continue
			}
			metricsBackend.GenerateMetrics(jobPosts)
		}
	}()
	return metricsBackend
}

// newPrometheusMetrics reads the status of the database, calculates all of the values of the metrics,
// registers these metrics to the prometheus registry, and returns the registry and the metrics.
// The registry and metrics can be modified from this struct from anywhere in the codebase.
func newPrometheusMetrics() *PrometheusMetrics {
	m := &PrometheusMetrics{
		registry: prometheus.NewRegistry(),
		JobPosts: jobPostsMetric(),

		RequestsTotal:    requestsTotalMetric(),
		RequestsDuration: requestDurationMetric(),
	}
	m.registry.MustRegister(m.JobPosts)

	m.registry.MustRegister(m.RequestsTotal)
	m.registry.MustRegister(m.RequestsDuration)

	m.registry.MustRegister(collectors.NewGoCollector())
	m.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	return m
}

// GenerateMetrics receives the live list of job posts to calculate the most recent values for the metrics
// defined for prometheus
func (pm *PrometheusMetrics) GenerateMetrics(jobPosts []db.JobPost) {
	pm.JobPosts.Set(float64(len(jobPosts)))
}

func jobPostsMetric() prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "job_posts_total",
		Help: "Total number of job posts",
	})
	return metric
}

func requestsTotalMetric() prometheus.CounterVec {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"method", "code"},
	)
	return *metric
}

func requestDurationMetric() prometheus.HistogramVec {
	metric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
		}, []string{"method", "code"},
	)
	return *metric
}
