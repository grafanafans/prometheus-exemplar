package cache

import (
	"context"

	"github.com/songjiayang/exemplar-demo/pkg/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Cache interface {
	Get(string, context.Context) interface{}
	Set(string, interface{}) error
}

func getWithOpel(ctx context.Context, spanName, key string, getFunc func() (bool, interface{})) interface{} {
	var (
		hit bool
		ret interface{}
	)

	_, span := otel.Tracer().Start(ctx, spanName)
	defer func() {
		span.SetAttributes(attribute.String("key", key), attribute.Bool("hit", hit))
		span.End()
	}()

	hit, ret = getFunc()

	return ret
}
