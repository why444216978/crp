package crp

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrCapacity = errors.New("capacity is 0")
)

type CacheReplacement interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Print() []interface{}
}

// TODO  调用用户的方法，需要 defer recover panic
