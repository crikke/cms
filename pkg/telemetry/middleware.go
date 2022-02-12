package telemetry

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.Register(httpReqCount)
}

const namespace = "service"

var (
	labels = []string{"code"}

	// httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	// 	Name: "http_response_time_seconds",
	// 	Help: "Duration of HTTP requests.",
	// }, []string{"path"})

	httpReqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests",
		}, labels)
)

func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		lvs := []string{fmt.Sprintf("%d", c.Writer.Status())}
		httpReqCount.WithLabelValues(lvs...).Inc()
	}
}
