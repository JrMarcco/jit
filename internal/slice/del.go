package slice

import "github.com/JrMarcco/easy-kit/internal/errs"

func Del[T any](slice []T, index int) ([]T, error) {
	length := len(slice)

	if index < 0 || index >= length {
		return nil, errs.IndexOutOfBoundsErr(length, index)
	}

	for i := index; i < length-1; i++ {
		slice[i] = slice[i+1]
	}

	return slice[:length-1], nil
}
