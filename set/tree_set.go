package set

import (
	easykit "github.com/JrMarcco/easy-kit"
	"github.com/JrMarcco/easy-kit/mapext"
)

var _ Set[any] = (*TreeSet[any])(nil)

type TreeSet[T any] struct {
	tm *mapext.TreeMap[T, struct{}]
}

func NewTreeSet[T any](cmp easykit.Comparator[T]) (*TreeSet[T], error) {
	tm, err := mapext.NewTreeMap[T, struct{}](cmp)
	if err != nil {
		return nil, err
	}

	return &TreeSet[T]{tm: tm}, nil
}

func (s *TreeSet[T]) Add(elem T) {
	_ = s.tm.Put(elem, struct{}{})
}

func (s *TreeSet[T]) Del(elem T) {
	_, _ = s.tm.Del(elem)
}

func (s *TreeSet[T]) Exist(elem T) bool {
	_, ok := s.tm.Get(elem)
	return ok
}

func (s *TreeSet[T]) Elems() []T {
	return s.tm.Keys()
}
