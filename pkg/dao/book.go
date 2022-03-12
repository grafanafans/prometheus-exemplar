package dao

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Book struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookService interface {
	List(int, int, context.Context) ([]*Book, error)
	Show(string, context.Context) (*Book, error)
}

type MockBookService struct{}

func NewMockBookService() *MockBookService {
	return &MockBookService{}
}
func (*MockBookService) List(page, perPage int, ctx context.Context) ([]*Book, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan := parentSpan.Tracer().StartSpan(
			"MockBookService.List",
			opentracing.ChildOf(parentSpan.Context()))
		defer childSpan.Finish()
	}

	// 5% with 200ms sleep and return nil
	if rand.Intn(100) <= 5 {
		time.Sleep(250 * time.Millisecond)
	}

	items := make([]*Book, 0, 10)

	for i := 1; i <= 10; i++ {
		items = append(items, &Book{
			Id:        fmt.Sprintf("book#%d", i),
			Name:      fmt.Sprintf("book-with-name-#%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return items, nil
}

func (*MockBookService) Show(id string, ctx context.Context) (*Book, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan := parentSpan.Tracer().StartSpan(
			"MockBookService.Show",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "id", Value: id})
		defer childSpan.Finish()
	}

	// 5% with 200ms sleep and return nil
	if rand.Intn(100) <= 5 {
		time.Sleep(250 * time.Millisecond)
	}

	return &Book{
		Id:        id,
		Name:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
