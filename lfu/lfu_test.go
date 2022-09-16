package lfu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	key       string
	value     string
	frequency int
}

var _ Data = (*data)(nil)

func (d *data) Key() string {
	return d.key
}

func (d *data) Value() interface{} {
	return d.value
}

func (d *data) Frequency() int {
	return d.frequency
}

func (d *data) SetFrequency(f int) {
	d.frequency = f
}

func (d *data) SetValue(v interface{}) {
	d.value = v.(string)
}

func New(key string, value interface{}) Data {
	return &data{key: key, value: value.(string), frequency: 1}
}

func TestLFU(t *testing.T) {
	o, _ := NewLFU(3, New)
	err := o.Put("1", "1")
	assert.Nil(t, err)
	err = o.Put("2", "2")
	assert.Nil(t, err)
	err = o.Put("3", "3")
	assert.Nil(t, err)

	dd := o.Print()
	assert.Equal(t, 3, len(dd))
	d := dd[0].(*data)
	assert.Equal(t, "3", d.Key())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 1, d.Frequency())
	d = dd[1].(*data)
	assert.Equal(t, "2", d.Key())
	assert.Equal(t, "2", d.Value())
	assert.Equal(t, 1, d.Frequency())
	d = dd[2].(*data)
	assert.Equal(t, "1", d.Key())
	assert.Equal(t, "1", d.Value())
	assert.Equal(t, 1, d.Frequency())

	_, err = o.Get("1")
	assert.Nil(t, err)
	dd = o.Print()
	d = dd[0].(*data)
	assert.Equal(t, "3", d.Key())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 1, d.Frequency())
	d = dd[1].(*data)
	assert.Equal(t, "2", d.Key())
	assert.Equal(t, "2", d.Value())
	assert.Equal(t, 1, d.Frequency())
	d = dd[2].(*data)
	assert.Equal(t, "1", d.Key())
	assert.Equal(t, "1", d.Value())
	assert.Equal(t, 2, d.Frequency())

	_, err = o.Get("2")
	assert.Nil(t, err)
	_, err = o.Get("3")
	assert.Nil(t, err)
	dd = o.Print()
	d = dd[0].(*data)
	assert.Equal(t, "3", d.Key())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 2, d.Frequency())
	d = dd[1].(*data)
	assert.Equal(t, "2", d.Key())
	assert.Equal(t, "2", d.Value())
	assert.Equal(t, 2, d.Frequency())
	d = dd[2].(*data)
	assert.Equal(t, "1", d.Key())
	assert.Equal(t, "1", d.Value())
	assert.Equal(t, 2, d.Frequency())

	err = o.Put("4", "4")
	assert.Nil(t, err)
	dd = o.Print()
	d = dd[0].(*data)
	assert.Equal(t, "4", d.Key())
	assert.Equal(t, "4", d.Value())
	assert.Equal(t, 1, d.Frequency())
	d = dd[1].(*data)
	assert.Equal(t, "3", d.Key())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 2, d.Frequency())
	d = dd[2].(*data)
	assert.Equal(t, "2", d.Key())
	assert.Equal(t, "2", d.Value())
	assert.Equal(t, 2, d.Frequency())
}
