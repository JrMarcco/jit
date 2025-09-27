package xmap

import (
	"errors"

	"github.com/JrMarcco/jit"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/JrMarcco/jit/internal/tree"
)

var _ imap[any, any] = (*TreeMap[any, any])(nil)

// TreeMap is a map implemented using a red-black tree.
type TreeMap[K any, V any] struct {
	tree *tree.RBTree[K, V]
}

func (tm *TreeMap[K, V]) Size() int64 {
	return tm.tree.Size()
}

func (tm *TreeMap[K, V]) Keys() []K {
	return tm.tree.Keys()
}

func (tm *TreeMap[K, V]) Vals() []V {
	return tm.tree.Vals()
}

func (tm *TreeMap[K, V]) Put(key K, val V) error {
	err := tm.tree.Put(key, val)
	if err != nil && errors.Is(err, errs.ErrSameRBNode) {
		// if the key already exists, update the value
		return tm.tree.Set(key, val)
	}
	return err
}

func (tm *TreeMap[K, V]) Del(key K) (V, bool) {
	v, err := tm.tree.Del(key)
	return v, err == nil
}

func (tm *TreeMap[K, V]) Get(key K) (V, bool) {
	v, err := tm.tree.Get(key)
	return v, err == nil
}

func (tm *TreeMap[K, V]) Iter(visitFunc func(key K, val V) bool) {
	tm.tree.Iter(visitFunc)
}

func (tm *TreeMap[K, V]) KeyVals() ([]K, []V) {
	return tm.tree.Kvs()
}

func NewTreeMap[K any, V any](cmp jit.Comparator[K]) (*TreeMap[K, V], error) {
	if cmp == nil {
		return nil, ErrNilComparator
	}
	return &TreeMap[K, V]{tree: tree.NewRBTree[K, V](cmp)}, nil
}

func NewTreeMapWithMap[K comparable, V any](cmp jit.Comparator[K], m map[K]V) (*TreeMap[K, V], error) {
	treeMap, err := NewTreeMap[K, V](cmp)
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		if err := treeMap.Put(k, v); err != nil {
			return nil, err
		}
	}

	return treeMap, nil
}
