package metrics

import (
	"demo/internal/common/config"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	DependencyDatabase   = "database"
	DependencyClickhouse = "clickhouse"
	DependencyRedis      = "redis"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

//go:generate mockery --name=IMetrics
type IMetrics interface {
	RecordHTTPRequest(method string, path string, statusCode int, duration time.Duration)
	RecordExternalCall(dependency string, operation string, status string, duration time.Duration)
	RecordNegativeKv(skId int, value float64)
}

type Metrics struct {
	// HTTP запросы
	HTTPRequestsTotal          *prometheus.CounterVec
	HTTPRequestDurationSeconds *prometheus.HistogramVec

	// Внешние зависимости
	ExternalCallsTotal           *prometheus.CounterVec
	ExternalCallsDurationSeconds *prometheus.HistogramVec

	// Бизнес метрики
	NegativeKvFinuslugiTotal *prometheus.CounterVec
}

func NewMetricsAndRegister(registry *prometheus.Registry, cfg *config.Config) IMetrics {
	metrics := NewMetrics()

	reg := prometheus.WrapRegistererWithPrefix(cfg.Metrics.Prefix, registry)

	reg.MustRegister(metrics.HTTPRequestsTotal)
	reg.MustRegister(metrics.HTTPRequestDurationSeconds)
	reg.MustRegister(metrics.ExternalCallsTotal)
	reg.MustRegister(metrics.ExternalCallsDurationSeconds)
	reg.MustRegister(metrics.NegativeKvFinuslugiTotal)

	return metrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP запросы
		HTTPRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество входящих HTTP запросов",
		}, []string{"method", "path", "status_class"}),

		HTTPRequestDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время обработки запросов в секундах",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status_class"}),

		// Внешние зависимости
		ExternalCallsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "external_calls_total",
			Help: "Количество запросов к внешним зависимостям",
		}, []string{"dependency", "operation", "status"}),

		ExternalCallsDurationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "external_calls_duration_seconds",
			Help:    "Время выполнения запросов к внешним сервисам в секундах",
			Buckets: prometheus.DefBuckets,
		}, []string{"dependency", "operation", "status"}),

		// Бизнес метрики
		NegativeKvFinuslugiTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "negative_kv_finuslugi_total",
			Help: "Отрицательное КВ Финуслуг",
		}, []string{"sk_id"}),
	}
}

func (m *Metrics) RecordHTTPRequest(method string, path string, statusCode int, duration time.Duration) {
	statusClass := strconv.Itoa(statusCode/100) + "xx"
	m.HTTPRequestsTotal.WithLabelValues(method, path, statusClass).Inc()
	m.HTTPRequestDurationSeconds.WithLabelValues(method, path, statusClass).Observe(duration.Seconds())
}

func (m *Metrics) RecordExternalCall(dependency string, operation string, status string, duration time.Duration) {
	m.ExternalCallsTotal.WithLabelValues(dependency, operation, status).Inc()
	m.ExternalCallsDurationSeconds.WithLabelValues(dependency, operation, status).Observe(duration.Seconds())
}

func (m *Metrics) RecordNegativeKv(skId int, value float64) {
	m.NegativeKvFinuslugiTotal.WithLabelValues(strconv.Itoa(skId)).Add(value)
}
