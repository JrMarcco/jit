package errs

import "fmt"

func NilErr(name string) error {
	return fmt.Errorf("[easy-kit] %s is nil", name)
}

func IndexOutOfBoundsErr(length int, index int) error {
	return fmt.Errorf("[easy-kit] index %d out of bounds for length %d", index, length)
}

func InvalidTypeErr(want string, got any) error {
	return fmt.Errorf("[easy-kit] invalid type: want %s, got %#v", want, got)
}

func EmptySliceErr() error {
	return fmt.Errorf("[easy-kit] slice is empty")
}

func InvalidKeyValLenErr() error {
	return fmt.Errorf("[easy-kit] keys and vals have different lengths")
}
