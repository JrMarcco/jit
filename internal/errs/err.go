package errs

import "fmt"

func IndexOutOfBoundsErr(length int, index int) error {
	return fmt.Errorf("[easy-kit] index %d out of bounds for length %d", index, length)
}

func InvalidTypeErr(want string, got any) error {
	return fmt.Errorf("[easy-kit] invalid type: want %s, got %#v", want, got)
}

func EmptySliceErr() error {
	return fmt.Errorf("[easy-kit] slice is empty")
}
