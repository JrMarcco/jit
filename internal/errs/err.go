package errs

import "fmt"

func IndexOutOfBoundsErr(len int, index int) error {
	return fmt.Errorf("[easy-kit] index %d out of bounds for length %d", index, len)
}

func InvalidTypeErr(want string, got any) error {
	return fmt.Errorf("[easy-kit] invalid type: want %s, got %#v", want, got)
}
