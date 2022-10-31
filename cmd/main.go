package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/songjiayang/exemplar-demo/pkg/api"
	"github.com/songjiayang/exemplar-demo/pkg/cache"
	"github.com/songjiayang/exemplar-demo/pkg/dao"
	"github.com/songjiayang/exemplar-demo/pkg/middleware"
	"github.com/songjiayang/exemplar-demo/pkg/otel"
)

var (
	appName      = "exemplar-demo"
	metricPath   = "/metrics"
	otelEndpoint = "collector:4318"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"/var/log/app.log"}
	logger, _ := cfg.Build()

	//set otel provider
	err := otel.SetTracerProvider(appName, "test", otelEndpoint)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// init mysql with retry
	go func() {
		retryCount := 10

		for i := 1; i <= retryCount; i++ {
			err := dao.InitMysqlDB()
			if err == nil {
				return
			}

			logger.Error(err.Error())
			time.Sleep(3 * time.Second)
		}

		log.Fatal("init mysql driver failed")
	}()

	r := gin.New()
	urlMapping := func(path string) string {
		switch path {
		case "/v1/books", "/ping":
			return path
		default:
			return "/v1/books/show"
		}
	}
	r.Use(middleware.Otel(metricPath, urlMapping))
	r.Use(middleware.Metrics(metricPath, urlMapping))

	myApi := api.NewApi(logger, cache.NewRedisCache(logger), dao.NewMysqlBookService())
	r.GET("/v1/books", myApi.Book.Index)
	r.GET("/v1/books/:id", myApi.Book.Show)

	// register prometheus metrics router
	proemtheusHandler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)
	r.GET(metricPath, func(ctx *gin.Context) {
		proemtheusHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.Run(":8080")
}
