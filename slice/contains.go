package slice

// ContainsFunc checks if the slice contains an element that satisfies the given function.
func ContainsFunc[T comparable](slice []T, elem T, eq equalFunc[T]) bool {
	for _, e := range slice {
		if eq(e, elem) {
			return true
		}
	}
	return false
}

// Contains checks if the slice contains the given element.
func Contains[T comparable](slice []T, elem T) bool {
	return ContainsFunc(slice, elem, func(src, dst T) bool { return src == dst })
}

// ContainsAnyFunc checks if the slice contains any of the given elements.
func ContainsAnyFunc[T comparable](slice []T, elems []T, eq equalFunc[T]) bool {
	for _, e := range elems {
		if ContainsFunc(slice, e, eq) {
			return true
		}
	}
	return false
}

// ContainsAny checks if the slice contains any of the given elements.
func ContainsAny[T comparable](slice []T, elems []T) bool {
	for _, e := range elems {
		if ContainsFunc(slice, e, func(src, dst T) bool { return src == dst }) {
			return true
		}
	}
	return false
}

// ContainsAllFunc checks if the slice contains all of the given elements.
func ContainsAllFunc[T comparable](slice []T, elems []T, eq equalFunc[T]) bool {
	if slice == nil || elems == nil {
		return false
	}

	for _, e := range elems {
		if !ContainsFunc(slice, e, eq) {
			return false
		}
	}
	return true
}

// ContainsAll checks if the slice contains all of the given elements.
func ContainsAll[T comparable](slice []T, elems []T) bool {
	if slice == nil || elems == nil {
		return false
	}

	for _, e := range elems {
		if !ContainsFunc(slice, e, func(src, dst T) bool { return src == dst }) {
			return false
		}
	}
	return true
}
