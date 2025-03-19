package slice

// Reverse reverses a slice in place.
func Reverse[T any](slice []T) []T {
	length := len(slice)
	reversed := make([]T, length)
	for i, v := range slice {
		reversed[length-i-1] = v
	}
	return reversed
}

// ReverseInPlace reverses a slice in place.
func ReverseInPlace[T any](slice []T) {
	length := len(slice)
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
