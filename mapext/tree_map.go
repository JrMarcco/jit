package mapext

import (
	"errors"

	easykit "github.com/JrMarcco/easy-kit"
	"github.com/JrMarcco/easy-kit/internal/errs"
	"github.com/JrMarcco/easy-kit/internal/tree"
)

var _ Map[any, any] = (*TreeMap[any, any])(nil)

// TreeMap is a map that is implemented using a red-black tree.
type TreeMap[K any, V any] struct {
	tree *tree.RBTree[K, V]
}

func NewTreeMap[K any, V any](cmp easykit.Comparator[K]) (*TreeMap[K, V], error) {
	if cmp == nil {
		return nil, ErrNilComparator
	}

	return &TreeMap[K, V]{tree: tree.NewRBTree[K, V](cmp)}, nil
}

func NewTreeMapWithMap[K comparable, V any](cmp easykit.Comparator[K], m map[K]V) (*TreeMap[K, V], error) {
	treeMap, err := NewTreeMap[K, V](cmp)
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		treeMap.Put(k, v)
	}

	return treeMap, nil
}

func (tm *TreeMap[K, V]) Put(k K, v V) error {
	err := tm.tree.Put(k, v)
	if err != nil && errors.Is(err, errs.ErrSameRBNode) {
		// if the key already exists, update the value
		return tm.tree.Set(k, v)
	}
	return err
}

func (tm *TreeMap[K, V]) Del(k K) (V, bool) {
	v, err := tm.tree.Del(k)
	return v, err == nil
}

func (tm *TreeMap[K, V]) Get(k K) (V, bool) {
	v, err := tm.tree.Get(k)
	return v, err == nil
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

func (tm *TreeMap[K, V]) KeyVals() ([]K, []V) {
	return tm.tree.Kvs()
}
