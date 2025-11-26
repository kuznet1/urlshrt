package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type myStruct struct {
	X int
}

func (m *myStruct) Reset() {
	m.X = 0
}

func TestPool(t *testing.T) {
	p := NewPool(func() *myStruct { return &myStruct{} })
	obj1 := p.Get()
	obj1.X = 42
	p.Put(obj1)
	obj2 := p.Get()
	assert.Equal(t, 0, obj2.X)
	assert.Equal(t, true, obj1 == obj2)
}
