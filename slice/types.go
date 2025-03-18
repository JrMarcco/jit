package slice

// equalFunc is a function that compares two elements and returns true if they are equal.
type equalFunc[T comparable] func(src, dst T) bool
