package jit

func Ptr[T any](v T) *T {
	return &v
}

func DePtr[T any](p *T) T {
	return *p
}
