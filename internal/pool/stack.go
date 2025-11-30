package pool

import "sync"

type Stack[T any] struct {
	mu   sync.Mutex
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data) == 0 {
		var zero T
		return zero, false
	}

	idx := len(s.data) - 1
	v := s.data[idx]
	s.data = s.data[:idx]
	return v, true
}
