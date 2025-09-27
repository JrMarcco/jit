package xslice

// Find returns the first element that matches the condition.
func Find[T any](slice []T, match matchFunc[T]) (T, bool) {
	for _, v := range slice {
		if match(v) {
			return v, true
		}
	}

	var t T
	return t, false
}

// FindAll returns all elements that match the condition.
func FindAll[T any](slice []T, match matchFunc[T]) []T {
	// estimate the capacity of the result slice
	// 25% of the slice length plus one
	res := make([]T, 0, len(slice)>>2+1)
	for _, v := range slice {
		if match(v) {
			res = append(res, v)
		}
	}
	return res
}
