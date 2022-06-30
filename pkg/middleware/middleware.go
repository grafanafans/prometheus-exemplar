package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/songjiayang/exemplar-demo/pkg/api"
	"go.opentelemetry.io/otel/attribute"

	"github.com/songjiayang/exemplar-demo/pkg/otel"
)

func Metrics(metricPath string, urlMapping func(string) string) gin.HandlerFunc {
	httpDurationsHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_durations_histogram_seconds",
		Help:    "Http latency distributions.",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2},
	}, []string{"method", "path", "code"})

	prometheus.MustRegister(httpDurationsHistogram)

	return func(c *gin.Context) {
		if c.Request.URL.Path == metricPath {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		url := urlMapping(c.Request.URL.Path)

		elapsed := float64(time.Since(start)) / float64(time.Second)
		observer := httpDurationsHistogram.WithLabelValues(method, url, status)
		observer.Observe(elapsed)

		if elapsed > 0.2 {
			observer.(prometheus.ExemplarObserver).ObserveWithExemplar(elapsed, prometheus.Labels{
				"traceID": c.GetHeader(api.XRequestID),
			})
		}
	}
}

func Otel(metricPath string, pathMapping func(string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == metricPath {
			c.Next()
			return
		}

		ctx, span := otel.Tracer().Start(c.Request.Context(), "root")
		defer span.End()

		span.SetAttributes(attribute.String("path", pathMapping(c.Request.URL.Path)))

		reqId := span.SpanContext().TraceID().String()
		c.Request.Header.Add(api.XRequestID, reqId)
		c.Header(api.XRequestID, reqId)

		c.Set("ctx", ctx)
		c.Next()
	}
}
