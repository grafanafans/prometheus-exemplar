package main

import (
	"context"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/songjiayang/exemplar-demo/pkg/api"
	"github.com/songjiayang/exemplar-demo/pkg/cache"
	"github.com/songjiayang/exemplar-demo/pkg/dao"
	"github.com/songjiayang/exemplar-demo/pkg/lokicore"
	"github.com/songjiayang/exemplar-demo/pkg/middleware"
	"github.com/songjiayang/exemplar-demo/pkg/otel"
)

var (
	appName    = "exemplar-demo"
	metricPath = "/metrics"
)

func main() {
	logger := newLokiLogger()

	//set otel provider
	tp, err := otel.SetTracerProvider(appName, "test")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer tp.Shutdown(context.Background())

	if err := dao.InitDB(); err != nil {
		logger.Fatal(err.Error())
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	urlMapping := newUrlMapping()
	r.Use(middleware.Otel(metricPath, urlMapping))
	r.Use(middleware.Metrics(metricPath, urlMapping))

	myApi := api.NewApi(logger, cache.NewRedisCache(logger), dao.NewMysqlBookService())
	r.GET("/v1/books", myApi.Book.Index)
	r.GET("/v1/books/:id", myApi.Book.Show)

	// register prometheus metrics router
	metricHandler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)

	r.GET(metricPath, func(ctx *gin.Context) {
		metricHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.Run(":8080")
}

func newLokiLogger() *zap.Logger {
	logger, _ := zap.NewProduction()

	cfg := &lokicore.LokiClientConfig{
		URL:       "http://loki:3100/api/prom/push",
		SendLevel: zapcore.InfoLevel,
		Labels: map[string]string{
			"app": appName,
		},
		TenantID: "demo",
	}

	lokiCore, err := lokicore.NewLokiCore(cfg)
	if err != nil {
		panic(err)
	}

	return logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, lokiCore)
	}))
}

func newUrlMapping() func(string) string {
	return func(path string) string {
		switch path {
		case "/v1/books", "/ping":
			return path
		default:
			return "/v1/books/show"
		}
	}
}
