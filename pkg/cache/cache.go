package cache

import "context"

type Cache interface {
	Get(string, context.Context) interface{}
	Set(string, interface{}) error
}
