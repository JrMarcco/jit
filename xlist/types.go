package xlist

type List[T any] interface {
	Insert(index int, val T) error
	Append(vals ...T) error
	Del(index int) error
	Set(index int, val T) error
	Get(index int) (T, error)
	Iter(visitFunc func(idx int, val T) error) error
	ToSlice() []T
	Cap() int
	Len() int
}
