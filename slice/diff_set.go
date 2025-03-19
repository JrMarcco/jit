package slice

// DiffSet returns the elements in src that are not in dst.
func DiffSet[T comparable](src, dst []T) []T {
	srcMap := toMap(src)
	for _, v := range dst {
		// remove v from srcMap
		delete(srcMap, v)
	}

	ret := make([]T, 0, len(srcMap))
	for k := range srcMap {
		ret = append(ret, k)
	}

	return ret
}

func DiffSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	slice := make([]T, 0, len(src))
	for _, v := range src {
		// check if v is in dst
		if !ContainsFunc(dst, v, func(t T) bool { return eq(v, t) }) {
			slice = append(slice, v)
		}
	}

	ret := make([]T, 0, len(slice))
	// remove duplicates
	for idx, val := range slice {
		if !ContainsFunc(slice[idx+1:], val, func(t T) bool { return eq(val, t) }) {
			ret = append(ret, val)
		}
	}

	return ret
}
