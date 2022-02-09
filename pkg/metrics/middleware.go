package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// TODO: better way to do this.
func WithHttpDuration(labels ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(labels...))
		c.Next()
		timer.ObserveDuration()
	}
}
