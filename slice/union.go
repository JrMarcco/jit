package slice

// UnionSetFunc returns the union of two sets.
func UnionSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(src)+len(dst))

	res = append(res, src...)
	res = append(res, dst...)

	return deDuplicateFunc(res, eq)
}

// UnionSet returns the union of two sets.
func UnionSet[T comparable](src, dst []T) []T {
	srcMap, dstMap := toMap(src), toMap(dst)
	for k := range srcMap {
		dstMap[k] = struct{}{}
	}

	ret := make([]T, 0, len(dstMap))
	for k := range dstMap {
		ret = append(ret, k)
	}

	return ret
}
