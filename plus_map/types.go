package plus_map

type IMap[K any, V any] interface {
	Size() int64
	Keys() []K
	Vals() []V
	Put(key K, val V) error
	Get(key K) (V, bool)
	Del(key K) (V, bool)
}
