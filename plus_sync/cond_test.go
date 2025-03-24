package plus_sync

import (
	"context"
	"sync"
	"testing"
)

func benchmarkCond(b *testing.B, waiterCnt int) {
	c := NewCond(&sync.Mutex{})
	done := make(chan bool)
	id := 0

	for routine := 0; routine < waiterCnt+1; routine++ {
		go func() {
			for i := 0; i < b.N; i++ {
				c.Locker.Lock()

				if id == -1 {
					c.Locker.Unlock()
					break
				}

				id++

				if id == waiterCnt+1 {
					id = 0
					c.Broadcast()
				} else {
					_ = c.Wait(context.Background())
				}

				c.Locker.Unlock()
			}

			c.Locker.Lock()
			id = -1
			c.Broadcast()
			c.Locker.Unlock()

			done <- true
		}()
	}
	for routine := 0; routine < waiterCnt+1; routine++ {
		<-done
	}

}

func BenchmarkCond_1(b *testing.B) {
	benchmarkCond(b, 1)
}

func BenchmarkCond_2(b *testing.B) {
	benchmarkCond(b, 2)
}

func BenchmarkCond_4(b *testing.B) {
	benchmarkCond(b, 4)
}

func BenchmarkCond_8(b *testing.B) {
	benchmarkCond(b, 8)
}

func BenchmarkCond_16(b *testing.B) {
	benchmarkCond(b, 16)
}

func BenchmarkCond_32(b *testing.B) {
	benchmarkCond(b, 32)
}

func BenchmarkCond_64(b *testing.B) {
	benchmarkCond(b, 64)
}

func BenchmarkCond_128(b *testing.B) {
	benchmarkCond(b, 128)
}
