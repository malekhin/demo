package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPMiddleware gin.HandlerFunc

func NewHTTPMiddleware(metrics IMetrics) HTTPMiddleware {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		status := c.Writer.Status()
		duration := time.Since(start)

		metrics.RecordHTTPRequest(method, path, status, duration)
	}
}
