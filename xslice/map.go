package xslice

// Map maps a slice to a new slice using a function.
func Map[Src any, Dst any](src []Src, fn func(idx int, src Src) Dst) []Dst {
	dst := make([]Dst, len(src))

	for i, v := range src {
		dst[i] = fn(i, v)
	}

	return dst
}

// FilterMap filters and maps a slice using a function.
func FilterMap[Src any, Dst any](src []Src, filter func(idx int, src Src) (Dst, bool)) []Dst {
	dst := make([]Dst, 0, len(src))

	for i, v := range src {
		if d, ok := filter(i, v); ok {
			dst = append(dst, d)
		}
	}

	return dst
}

// ToMap converts a slice to a map.
// the key is the result of the function.
func ToMap[K comparable, V any](slice []V, fn func(elem V) K) map[K]V {
	return ToMapWithVal(slice, func(elem V) (K, V) { return fn(elem), elem })
}

// ToMapWithVal converts a slice to a map.
// the key and value are the result of the function.
func ToMapWithVal[E any, K comparable, V any](slice []E, fn func(elem E) (K, V)) map[K]V {
	res := make(map[K]V, len(slice))
	for _, v := range slice {
		k, v := fn(v)
		res[k] = v
	}
	return res
}

// toMap converts a slice to a map.
func toMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range slice {
		// use empty struct to save memory
		m[v] = struct{}{}
	}
	return m
}

// deDuplicateFunc returns a slice of unique elements.
func deDuplicateFunc[T comparable](slice []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(slice))
	for i, v := range slice {
		if !ContainsFunc(slice[i+1:], func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}
	return res
}

// deDuplicate returns a slice of unique elements.
func deDuplicate[T comparable](slice []T) []T {
	m := toMap(slice)
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
