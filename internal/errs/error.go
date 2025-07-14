package errs

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrSameRBNode   = errors.New("[easy-kit] cannot insert same red-black tree node")
	ErrNodeNotFound = errors.New("[easy-kit] cannot find node in red-black tree")
)

func NilErr(name string) error {
	return fmt.Errorf("[easy-kit] %s is nil", name)
}

func ErrIndexOutOfBounds(length int, index int) error {
	return fmt.Errorf("[easy-kit] index %d out of bounds for length %d", index, length)
}

func ErrEmptySlice() error {
	return fmt.Errorf("[easy-kit] slice is empty")
}

func ErrInvalidKeyValLen() error {
	return fmt.Errorf("[easy-kit] keys and vals have different lengths")
}

func ErrInvalidInterval(interval time.Duration) error {
	return fmt.Errorf("[easy-kit] invalid interval: %v, expected interval value should greater than 0", interval)
}

func ErrInvalidMaxInterval(maxInterval time.Duration) error {
	return fmt.Errorf(
		"[easy-kit] invalid max interval: %v, expected max interval value should greater than init interval",
		maxInterval,
	)
}

func ErrRetryTimeExhausted(latestErr error) error {
	return fmt.Errorf("[easy-kit] retry time exhausted, the latest error: %w", latestErr)
}
