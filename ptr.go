package easy_kit

func Ptr[T any](v T) *T {
	return &v
}

func DePtr[T any](p *T) T {
	return *p
}
