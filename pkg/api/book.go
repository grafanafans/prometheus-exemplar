package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/songjiayang/exemplar-demo/pkg/dao"
)

type Book struct {
	Base    *Api
	Service dao.BookService
}

func (b *Book) Index(ctx *gin.Context) {
	logger := b.Base.ContextLogger(ctx)
	spanCtx := ctx.MustGet("ctx").(context.Context)

	page, _ := strconv.ParseInt(ctx.Param("page"), 10, 64)
	perPage, _ := strconv.ParseInt(ctx.Param("perPage"), 10, 64)

	logger.Info("start list books with cache")
	if cache := b.Base.cache.Get("book:list", spanCtx); cache != nil {
		logger.Info("listed books with cache")
		ctx.JSON(200, cache)
		return
	}

	logger.Info("start list books with db service")
	books, err := b.Service.List(int(page), int(perPage), spanCtx)
	if err != nil {
		logger.Error(fmt.Sprintf("list books with error: %v", err))
		return
	}

	logger.Info("start setting books cache")
	b.Base.cache.Set("book:list", books)
	logger.Info("end setting books cache")

	ctx.JSON(200, books)
}

func (b *Book) Show(ctx *gin.Context) {
	logger := b.Base.ContextLogger(ctx)
	spanCtx := ctx.MustGet("ctx").(context.Context)

	id := ctx.Param("id")
	cacheKey := "book:show:" + id

	logger.Info(fmt.Sprintf("get book info with cache key %s", cacheKey))
	if cache := b.Base.cache.Get(cacheKey, spanCtx); cache != nil {
		ctx.JSON(200, cache)
		return
	}

	logger.Info("start get book info from db service")
	book, err := b.Service.Show(id, spanCtx)
	if err != nil {
		logger.Error(fmt.Sprintf("show book with error: %v", err))
		ctx.JSON(500, err.Error())
		return
	}

	logger.Info("start setting book cache")
	b.Base.cache.Set(cacheKey, book)
	logger.Info("end setting book cache")
	ctx.JSON(200, book)
}
