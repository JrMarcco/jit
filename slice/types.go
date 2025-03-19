package slice

// eqFunc is a function that compares two elements and returns true if they are equal.
type eqFunc[T comparable] func(src, dst T) bool
