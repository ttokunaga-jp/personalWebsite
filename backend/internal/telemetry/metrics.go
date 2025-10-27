package telemetry

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/takumi/personal-website/internal/config"
)

// Metrics instruments HTTP handlers with Prometheus metrics.
type Metrics struct {
	enabled         bool
	endpoint        string
	registry        *prometheus.Registry
	requestDuration *prometheus.HistogramVec
	requestTotal    *prometheus.CounterVec
	requestErrors   *prometheus.CounterVec
}

// NewMetrics constructs the collector and registers the metrics.
func NewMetrics(cfg *config.AppConfig) *Metrics {
	if cfg == nil {
		return &Metrics{enabled: false}
	}

	metricsCfg := cfg.Metrics
	if !metricsCfg.Enabled {
		return &Metrics{enabled: false}
	}

	registry := prometheus.NewRegistry()
	labels := []string{"method", "path"}

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsCfg.Namespace,
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests",
		Buckets:   prometheus.DefBuckets,
	}, labels)

	requestTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsCfg.Namespace,
		Name:      "http_requests_total",
		Help:      "Total number of processed HTTP requests",
	}, append(labels, "status"))

	requestErrors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsCfg.Namespace,
		Name:      "http_request_errors_total",
		Help:      "Number of HTTP requests resulting in errors (>=500)",
	}, labels)

	registry.MustRegister(requestDuration, requestTotal, requestErrors)

	endpoint := metricsCfg.Endpoint
	if endpoint == "" {
		endpoint = "/metrics"
	}

	return &Metrics{
		enabled:         true,
		endpoint:        endpoint,
		registry:        registry,
		requestDuration: requestDuration,
		requestTotal:    requestTotal,
		requestErrors:   requestErrors,
	}
}

// Handler instruments requests when metrics are enabled.
func (m *Metrics) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || !m.enabled {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		duration := time.Since(start).Seconds()
		labels := prometheus.Labels{"method": c.Request.Method, "path": path}
		m.requestDuration.With(labels).Observe(duration)

		statusCode := c.Writer.Status()
		m.requestTotal.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
			"status": strconv.Itoa(statusCode),
		}).Inc()

		if statusCode >= 500 {
			m.requestErrors.With(labels).Inc()
		}
	}
}

// Register exposes the Prometheus endpoint on the supplied router.
func (m *Metrics) Register(router *gin.Engine) {
	if m == nil || !m.enabled || router == nil {
		return
	}

	router.GET(m.endpoint, gin.WrapH(promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})))
}
