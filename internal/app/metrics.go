package app

import (
	"demo/internal/common/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

var MetricsModule = fx.Module("metrics",
	fx.Provide(metrics.NewMetricsAndRegister),
	fx.Provide(metrics.NewHTTPMiddleware),
	fx.Provide(NewPrometheusRegistry),
	fx.Provide(NewMetricsHandler),
)

type MetricsHandler gin.HandlerFunc

func NewMetricsHandler(registry *prometheus.Registry) MetricsHandler {
	return MetricsHandler(gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))
}

func NewPrometheusRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
