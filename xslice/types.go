package xslice

// eqFunc is a function that compares two elements and returns true if they are equal.
type eqFunc[T comparable] func(src, dst T) bool

// matchFunc is a function that returns true if the element matches the condition.
type matchFunc[T any] func(T) bool
