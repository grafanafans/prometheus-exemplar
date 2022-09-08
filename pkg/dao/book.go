package dao

import (
	"context"
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

type MysqlBookService struct{}

func NewMysqlBookService() *MysqlBookService {
	return &MysqlBookService{}
}
func (*MysqlBookService) List(page, perPage int, ctx context.Context) (items []*Book, err error) {
	_, span := otel.Tracer().Start(ctx, "MysqlBookService.List")
	defer span.End()

	// 10% with random sleep
	if rand.Intn(100) <= 10 {
		time.Sleep(time.Duration(rand.Intn(350)) * time.Millisecond)
	}

	err = db.Find(&items).Limit(10).Error
	return
}

func (*MysqlBookService) Show(id string, ctx context.Context) (item *Book, err error) {
	_, span := otel.Tracer().Start(ctx, "MysqlBookService.Show")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

	// 10% with random sleep
	if rand.Intn(100) <= 10 {
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	}

	err = db.Where(Book{Id: id}).Find(&item).Error
	return
}
