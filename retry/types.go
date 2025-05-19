package retry

import "time"

type Strategy interface {
	Next() (time.Duration, bool)
	Report(err error) Strategy
}
