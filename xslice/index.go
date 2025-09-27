package xslice

// IndexFunc returns the index of the first element that matches the condition.
func IndexFunc[T comparable](slice []T, match matchFunc[T]) int {
	for k, v := range slice {
		if match(v) {
			return k
		}
	}
	return -1
}

// Index returns the index of the first occurrence of elem in slice.
func Index[T comparable](slice []T, elem T) int {
	return IndexFunc(slice, func(t T) bool { return t == elem })
}

// LastIndexFunc returns the index of the last element that matches the condition.
func LastIndexFunc[T comparable](slice []T, match matchFunc[T]) int {
	for k := len(slice) - 1; k >= 0; k-- {
		if match(slice[k]) {
			return k
		}
	}
	return -1
}

// LastIndex returns the index of the last occurrence of elem in slice.
func LastIndex[T comparable](slice []T, elem T) int {
	return LastIndexFunc(slice, func(t T) bool { return t == elem })
}

// IndexAllFunc returns all indices of elements that match the condition.
func IndexAllFunc[T comparable](slice []T, match matchFunc[T]) []int {
	res := make([]int, 0, len(slice)>>2+1)
	for k, v := range slice {
		if match(v) {
			res = append(res, k)
		}
	}
	return res
}

// IndexAll returns all indices of elements that match the condition.
func IndexAll[T comparable](slice []T, elem T) []int {
	return IndexAllFunc(slice, func(t T) bool { return t == elem })
}
