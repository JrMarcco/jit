package xlist

import "sync"

var _ List[any] = (*ConcurrentList[any])(nil)

type ConcurrentList[T any] struct {
	List[T]
	mu sync.RWMutex
}

func (cl *ConcurrentList[T]) Insert(index int, val T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Insert(index, val)
}

func (cl *ConcurrentList[T]) Append(vals ...T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Append(vals...)
}

func (cl *ConcurrentList[T]) Del(index int) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Del(index)
}

func (cl *ConcurrentList[T]) Set(index int, val T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Set(index, val)
}

func (cl *ConcurrentList[T]) Get(index int) (T, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Get(index)
}

func (cl *ConcurrentList[T]) Iter(visitFunc func(idx int, val T) error) error {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Iter(visitFunc)
}

func (cl *ConcurrentList[T]) ToSlice() []T {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.ToSlice()
}

func (cl *ConcurrentList[T]) Cap() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Cap()
}

func (cl *ConcurrentList[T]) Len() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Len()
}
