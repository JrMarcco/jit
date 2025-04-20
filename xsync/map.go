package xsync

import "sync"

// Map is generic packing for sync.Map.
type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := m.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), ok
}

func (m *Map[K, V]) Store(key K, val V) {
	m.m.Store(key, val)
}

func (m *Map[K, V]) LoadOrStore(ket K, val V) (V, bool) {
	v, ok := m.m.LoadOrStore(ket, val)
	if v != nil {
		return v.(V), ok
	}

	var zero V
	return zero, ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	v, ok := m.m.LoadAndDelete(key)
	if v != nil {
		return v.(V), ok
	}

	var zero V
	return zero, ok
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Range(fn func(key K, val V) bool) {
	m.m.Range(func(key, val any) bool {
		var (
			k K
			v V
		)

		if val != nil {
			v = val.(V)
		}

		if key != nil {
			k = key.(K)
		}

		return fn(k, v)
	})
}
