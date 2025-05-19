package retry

import (
	"math/bits"
	"sync/atomic"
	"time"
)

var _ Strategy = (*AdaptiveTimeoutStrategy)(nil)

type AdaptiveTimeoutStrategy struct {
	strategy  Strategy // basic retry strategy
	threshold int      // timeout threshold

	bufferSize int      // size of the slide window
	ringBuffer []uint64 // using as a slide window to store timeout information

	totalBit uint64
}

func (a *AdaptiveTimeoutStrategy) Next() (time.Duration, bool) {
	failureCnt := a.getFailureCnt()
	if failureCnt >= a.threshold {
		return 0, false
	}
	return a.strategy.Next()
}

func (a *AdaptiveTimeoutStrategy) NextWithRetried(retriedTimes int32) (time.Duration, bool) {
	failureCnt := a.getFailureCnt()
	if failureCnt >= a.threshold {
		return 0, false
	}

	if s, ok := a.strategy.(interface {
		NextWithRetried(int32) (time.Duration, bool)
	}); ok {
		return s.NextWithRetried(retriedTimes)
	}

	// fallback: just call Next()
	return a.strategy.Next()
}

func (a *AdaptiveTimeoutStrategy) Report(err error) Strategy {
	if err == nil {
		a.markAsSuccess()
		return a
	}

	a.markAsFailure()
	return a
}

func (a *AdaptiveTimeoutStrategy) markAsSuccess() {}

func (a *AdaptiveTimeoutStrategy) markAsFailure() {}

func (a *AdaptiveTimeoutStrategy) getFailureCnt() int {
	var cnt int
	for i := 0; i < a.bufferSize; i++ {
		val := atomic.LoadUint64(&a.ringBuffer[i])
		cnt += bits.OnesCount64(val)
	}

	return cnt
}

func NewAdaptiveTimeoutStrategy(strategy Strategy, bufferSize int, threshold int) *AdaptiveTimeoutStrategy {
	return &AdaptiveTimeoutStrategy{
		strategy:   strategy,
		threshold:  threshold,
		bufferSize: bufferSize,
		ringBuffer: make([]uint64, bufferSize),
		totalBit:   uint64(64) & uint64(bufferSize),
	}
}
