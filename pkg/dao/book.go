package dao

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/songjiayang/exemplar-demo/pkg/otel"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Book struct {
	Id        string    `json:"id" gorm:"uniqueIndex"`
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
	_, span := otel.Tracer().Start(ctx, "MockBookService.List")
	defer span.End()

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
	_, span := otel.Tracer().Start(ctx, "MockBookService.Show")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

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

type MysqlBookService struct{}

func NewMysqlBookService() *MysqlBookService {
	return &MysqlBookService{}
}

func (*MysqlBookService) List(page, perPage int, ctx context.Context) (items []*Book, err error) {
	_, span := otel.Tracer().Start(ctx, "MysqlBookService.List")
	defer span.End()

	err = db.Find(&items).Limit(10).Error
	return
}

func (*MysqlBookService) Show(id string, ctx context.Context) (item *Book, err error) {
	_, span := otel.Tracer().Start(ctx, "MysqlBookService.Show")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

	err = db.Where(Book{Id: id}).Find(&item).Error
	return
}
