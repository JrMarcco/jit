package slice

// toMap converts a slice to a map.
func toMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range slice {
		// use empty struct to save memory
		m[v] = struct{}{}
	}
	return m
}
