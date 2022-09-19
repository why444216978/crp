package lfu

import (
	"container/list"
	"sort"

	"github.com/pkg/errors"
	wp "github.com/why444216978/go-util/panic"

	"github.com/why444216978/crp"
)

type Data interface {
	Key() string
	Value() interface{}
	SetValue(v interface{})
	Frequency() int
	SetFrequency(f int)
}

type lfu struct {
	keys           map[string]*list.Element // save key map
	frequencyLists map[int]*list.List       // save the same frequency list
	capacity       int
	min            int // save min frequency
	newFunc        func(string, interface{}) Data
}

var _ crp.CacheReplacement = (*lfu)(nil)

func NewLFU(capacity int, f func(string, interface{}) Data) (*lfu, error) {
	if capacity == 0 {
		return nil, crp.ErrCapacity
	}
	if f == nil {
		return nil, errors.New("newFunc nil")
	}
	return &lfu{
		keys:           map[string]*list.Element{},
		frequencyLists: map[int]*list.List{},
		capacity:       capacity,
		min:            0,
		newFunc:        f,
	}, nil
}

func (c *lfu) Get(key string) (interface{}, error) {
	el, ok := c.keys[key]
	if !ok {
		return nil, crp.ErrNotFound
	}

	// data := c.assertValue()
	currentNode, err := c.assertValue(el)
	if err != nil {
		return nil, err
	}

	oldFrequency := currentNode.Frequency()
	newFrequency := oldFrequency + 1
	currentNode.SetFrequency(newFrequency)

	// remove current node from old frequency list
	c.frequencyLists[oldFrequency].Remove(el)

	c.add(currentNode)

	// if old frequency equal min and empty, new frequency is assigned to min
	if oldFrequency == c.min && c.frequencyLists[oldFrequency].Len() == 0 {
		c.min = newFrequency
	}

	return currentNode.Value(), nil
}

func (c *lfu) Put(key string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wp.NewPanicError(r)
		}
	}()

	if c.capacity == 0 {
		return crp.ErrCapacity
	}

	// if exists, call Get function to incr visits number
	if currentValue, ok := c.keys[key]; ok {
		data, err := c.assertValue(currentValue)
		if err != nil {
			return err
		}
		data.SetValue(value)
		_, _ = c.Get(key)
		return nil
	}

	// if cache full
	// delete min list last node
	// delete map
	if c.capacity == len(c.keys) {
		minList := c.frequencyLists[c.min]
		lastNode := minList.Back()
		data, err := c.assertValue(lastNode)
		if err != nil {
			return err
		}
		delete(c.keys, data.Key())
		minList.Remove(lastNode)
	}

	// add new node
	minFrequency := 1
	c.min = minFrequency
	currentNode := c.newFunc(key, value)
	c.add(currentNode)

	return
}

func (c *lfu) add(data Data) {
	f := data.Frequency()

	newList, ok := c.frequencyLists[f]
	if !ok {
		c.frequencyLists[f] = list.New()
		newList = c.frequencyLists[f]
	}

	// push node
	newNode := newList.PushFront(data)

	// save node map
	c.keys[data.Key()] = newNode
}

func (c *lfu) assertValue(el *list.Element) (Data, error) {
	data, ok := el.Value.(Data)
	if !ok {
		return nil, errors.Errorf("convert to Data fail source type:%T", el.Value)
	}
	return data, nil
}

func (c *lfu) Print() []interface{} {
	arr := []interface{}{}

	var keys []int
	for k := range c.frequencyLists {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		root := c.frequencyLists[k].Front()
		for {
			if root == nil {
				break
			}

			arr = append(arr, root.Value)
			root = root.Next()
		}
	}

	return arr
}
