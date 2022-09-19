package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/why444216978/crp"
)

type data struct {
	key   string
	value string
}

var _ Data = (*data)(nil)

func (d *data) Key() string {
	return d.key
}

func (d *data) Value() interface{} {
	return d.value
}

func (d *data) SetValue(v interface{}) {
	d.value = v.(string)
}

func New(key string, value interface{}) Data {
	return &data{key: key, value: value.(string)}
}

func TestLRU(t *testing.T) {
	o, err := NewLRU(3, New)
	assert.Nil(t, err)
	err = o.Put("a", "a")
	assert.Nil(t, err)
	err = o.Put("b", "b")
	assert.Nil(t, err)
	err = o.Put("c", "c")
	assert.Nil(t, err)

	dd := o.Print()
	assert.Equal(t, 3, len(dd))
	d := dd[0].(*data)
	assert.Equal(t, "c", d.Key())
	assert.Equal(t, "c", d.Value())
	d = dd[1].(*data)
	assert.Equal(t, "b", d.Key())
	assert.Equal(t, "b", d.Value())
	d = dd[2].(*data)
	assert.Equal(t, "a", d.Key())
	assert.Equal(t, "a", d.Value())

	err = o.Put("a", "a")
	assert.Nil(t, err)
	dd = o.Print()
	assert.Equal(t, 3, len(dd))
	d = dd[0].(*data)
	assert.Equal(t, "a", d.Key())
	assert.Equal(t, "a", d.Value())
	d = dd[1].(*data)
	assert.Equal(t, "c", d.Key())
	assert.Equal(t, "c", d.Value())
	d = dd[2].(*data)
	assert.Equal(t, "b", d.Key())
	assert.Equal(t, "b", d.Value())

	// get non-exists
	b, err := o.Get("non")
	assert.Equal(t, crp.ErrNotFound, err)
	assert.Nil(t, b)

	// get exists
	b, err = o.Get("b")
	assert.Nil(t, err)
	bb := b.(*data)
	assert.Equal(t, "b", bb.key)
	assert.Equal(t, "b", bb.value)
	dd = o.Print()
	d = dd[0].(*data)
	assert.Equal(t, "b", d.key)
	assert.Equal(t, "b", d.value)
	d = dd[1].(*data)
	assert.Equal(t, "a", d.key)
	assert.Equal(t, "a", d.value)
	d = dd[2].(*data)
	assert.Equal(t, "c", d.key)
	assert.Equal(t, "c", d.value)
}
