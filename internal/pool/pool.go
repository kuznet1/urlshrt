package pool

import "sync"

type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	mu      sync.Mutex
	items   []T
	newFunc func() T
}

func NewPool[T Resettable](newFunc func() T) *Pool[T] {
	return &Pool[T]{newFunc: newFunc}
}

func (p *Pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := len(p.items)
	if n == 0 {
		return p.newFunc()
	}

	obj := p.items[n-1]
	p.items = p.items[:n-1]
	return obj
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.mu.Lock()
	defer p.mu.Unlock()

	p.items = append(p.items, obj)
}
