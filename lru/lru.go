package lru

import (
	"container/list"
	"sync"

	"github.com/pkg/errors"
	wp "github.com/why444216978/go-util/panic"

	"github.com/why444216978/crp"
)

type Data interface {
	Key() string
	Value() interface{}
	SetValue(v interface{})
}

type lru struct {
	lock    sync.Mutex
	Cap     int
	Keys    map[string]*list.Element
	List    *list.List
	newFunc func(string, interface{}) Data
}

var _ crp.CacheReplacement = (*lru)(nil)

func NewLRU(capacity int, f func(string, interface{}) Data) (crp.CacheReplacement, error) {
	if f == nil {
		return nil, errors.New("newFunc nill")
	}

	return &lru{
		Cap:     capacity,
		Keys:    map[string]*list.Element{},
		List:    list.New(),
		newFunc: f,
	}, nil
}

func (c *lru) Get(key string) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if el, ok := c.Keys[key]; ok {
		c.List.MoveToFront(el)
		return c.assertValue(el)
	}
	return nil, crp.ErrNotFound
}

func (c *lru) Put(key string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wp.NewPanicError(r)
		}
	}()

	c.lock.Lock()
	defer c.lock.Unlock()

	if el, ok := c.Keys[key]; ok {
		data, err := c.assertValue(el)
		if err != nil {
			return err
		}
		data.SetValue(value)
		c.List.MoveToFront(el)
	} else {
		e := c.List.PushFront(c.newFunc(key, value))
		c.Keys[key] = e
	}

	if c.List.Len() > c.Cap {
		e := c.List.Back()
		c.List.Remove(e)
		delete(c.Keys, e.Value.(Data).Key())
	}

	return
}

func (c *lru) Print() []interface{} {
	arr := make([]interface{}, c.List.Len())

	i := 0
	root := c.List.Front()
	for {
		if root == nil {
			break
		}
		arr[i] = root.Value

		root = root.Next()
		i = i + 1
	}

	return arr
}

func (c *lru) assertValue(el *list.Element) (Data, error) {
	data, ok := el.Value.(Data)
	if !ok {
		return nil, errors.Errorf("convert to Data fail source type:%T", el.Value)
	}
	return data, nil
}
