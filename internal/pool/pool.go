package pool

type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	items   Stack[T]
	newFunc func() T
}

func NewPool[T Resettable](newFunc func() T) *Pool[T] {
	return &Pool[T]{newFunc: newFunc}
}

func (p *Pool[T]) Get() T {
	if v, ok := p.items.Pop(); ok {
		return v
	}
	return p.newFunc()
}

func (p *Pool[T]) Put(v T) {
	v.Reset()
	p.items.Push(v)
}
