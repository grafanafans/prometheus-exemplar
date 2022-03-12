package middleware

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/songjiayang/exemplar-demo/pkg/api"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
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

func Jaeger(serviceName, jaegerServer string, pathMapping func(string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		tracer, closer := newTracer(serviceName, jaegerServer)
		defer closer.Close()
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(pathMapping(c.Request.URL.Path))
		} else {
			parentSpan = opentracing.StartSpan(
				pathMapping(c.Request.URL.Path),
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
		}
		defer parentSpan.Finish()

		sc, ok := parentSpan.Context().(jaeger.SpanContext)
		if ok {
			reqId := sc.TraceID().String()
			c.Request.Header.Add(api.XRequestID, reqId)
			c.Header(api.XRequestID, reqId)
		}

		c.Set("tracer", tracer)
		c.Set("ctx", opentracing.ContextWithSpan(context.Background(), parentSpan))
		c.Next()
	}
}

func newTracer(service, host string) (opentracing.Tracer, io.Closer) {
	cfg := jaegerConfig.Configuration{
		ServiceName: service,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: host,
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
