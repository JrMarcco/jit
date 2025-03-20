package slice

// IntersectSetFunc returns the intersection of two slices.
func IntersectSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(src))
	for _, v := range src {
		if ContainsFunc(dst, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}
	return deDuplicateFunc(res, eq)
}

// IntersectSet returns the intersection of two slices.
func IntersectSet[T comparable](src, dst []T) []T {
	m := toMap(src)
	res := make([]T, 0, len(m))
	for _, v := range dst {
		if _, ok := m[v]; ok {
			res = append(res, v)
		}
	}
	return deDuplicate(res)
}
