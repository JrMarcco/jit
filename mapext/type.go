package mapext

type Map[K any, V any] interface {
	Size() int64
	Keys() []K
	Vals() []V
	Put(key K, val V) error
	Del(key K) (V, bool)
	Get(key K) (V, bool)
}
