package slice

// toMap converts a slice to a map.
func toMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range slice {
		// use empty struct to save memory
		m[v] = struct{}{}
	}
	return m
}

// duplicateFunc returns a slice of unique elements.
func duplicateFunc[T comparable](slice []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(slice))
	for i, v := range slice {
		if !ContainsFunc(slice[i+1:], func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}
	return res
}

// duplicate returns a slice of unique elements.
func duplicate[T comparable](slice []T) []T {
	m := toMap(slice)
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
