package slice

import (
	"github.com/JrMarcco/easy-kit/internal/errs"
)

func Add[T any](slice []T, index int, item T) ([]T, error) {
	length := len(slice)

	if index < 0 || index > length {
		return nil, errs.IndexOutOfBoundsErr(length, index)
	}

	if index == length {
		return append(slice, item), nil
	}

	// 扩容一个位置，length + 1
	var zeroVal T
	slice = append(slice, zeroVal)

	// 注意 length 是扩容后的长度所以不需要减 1
	for i := length; i > index; i-- {
		if i-1 >= 0 {
			slice[i] = slice[i-1]
		}
	}

	slice[index] = item

	return slice, nil
}
