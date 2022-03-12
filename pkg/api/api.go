package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/songjiayang/exemplar-demo/pkg/cache"
	"github.com/songjiayang/exemplar-demo/pkg/dao"
)

const XRequestID = "X-Request-ID"

type Api struct {
	logger *zap.Logger
	Book   *Book
	cache  cache.Cache
}

func NewApi(logger *zap.Logger, cache cache.Cache, bs dao.BookService) *Api {
	a := &Api{
		logger: logger,
		cache:  cache,
	}

	a.Book = &Book{a, bs}
	return a
}

func (a *Api) ContextLogger(ctx *gin.Context) *zap.Logger {
	return a.logger.With(zap.String("traceID", ctx.GetHeader(XRequestID)))
}
