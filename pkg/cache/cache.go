package cache

import (
	"context"
	"math/rand"
	"time"

	"github.com/songjiayang/exemplar-demo/pkg/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Cache interface {
	Get(string, context.Context) interface{}
	Set(string, interface{}) error
}

func getWithOtel(ctx context.Context, spanName, key string, getFunc func() (bool, interface{})) interface{} {
	var (
		hit bool
		ret interface{}
	)

	_, span := otel.Tracer().Start(ctx, spanName)
	defer func() {
		span.SetAttributes(attribute.String("key", key), attribute.Bool("hit", hit))
		span.End()
	}()

	// random sleep for mock
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	hit, ret = getFunc()

	return ret
}
