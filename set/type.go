package set

type Set[T comparable] interface {
	Add(key T)
	Del(key T)
	Exist(key T) bool
	Elems() []T
}
