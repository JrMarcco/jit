package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type shutdownNowRes struct {
	tasks []Task
	err   error
}

func poolInternalState(p *BlockTaskPool) int32 {
	for {
		state := atomic.LoadInt32(&p.state)
		if state != stateLocked {
			return state
		}
	}

}

func runningPool(t *testing.T, initG, queueSize int32) *BlockTaskPool {
	p, err := NewBlockTaskPool(initG, queueSize)
	assert.NoError(t, err)

	assert.Equal(t, poolInternalState(p), stateCreated)
	assert.NoError(t, p.Start())
	assert.Equal(t, poolInternalState(p), stateRunning)
	return p
}

func closedPool(t *testing.T, initG, queueSize int32) *BlockTaskPool {
	p := runningPool(t, initG, queueSize)

	tasks, err := p.ShutdownNow()
	assert.NoError(t, err)
	assert.Equal(t, len(tasks), 0)
	assert.Equal(t, poolInternalState(p), stateClosed)
	return p
}

func runningPoolWithFilledQueue(t *testing.T, initG, queueSize int32) (*BlockTaskPool, chan struct{}) {
	p := runningPool(t, initG, queueSize)
	wait := make(chan struct{})

	for i := int32(0); i < initG+queueSize; i++ {
		err := p.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			<-wait
			return nil
		}))
		assert.NoError(t, err)
	}
	return p, wait
}

func TestTaskPool_State(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		ch, err := p.State(context.Background(), time.Millisecond)
		assert.NoError(t, err)
		assert.NotZero(t, ch)

		done, err := p.Shutdown()
		assert.NoError(t, err)
		<-done
	})

	t.Run("call after context cancel", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = p.State(ctx, time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("call after shutdown", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		done, err := p.Shutdown()
		assert.NoError(t, err)

		<-done

		_, err = p.State(context.Background(), time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("call after shutdown now", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		_, err = p.ShutdownNow()
		assert.NoError(t, err)

		_, err = p.State(context.Background(), time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("close chan after context timeout", func(t *testing.T) {
		t.Parallel()

		initG, queueSize := int32(1), int32(3)
		p, waitCh := runningPoolWithFilledQueue(t, initG, queueSize)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
		stateCh, err := p.State(ctx, time.Millisecond)
		assert.NoError(t, err)

		go func() {
			// 模拟 context 超时
			<-time.After(3 * time.Millisecond)
			cancel()
		}()

		for {
			state, ok := <-stateCh
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}

		close(waitCh)
		_, err = p.Shutdown()
		assert.NoError(t, err)
	})

	t.Run("close chan after context canceled", func(t *testing.T) {
		t.Parallel()

		initG, queueSize := int32(1), int32(3)
		p, waitCh := runningPoolWithFilledQueue(t, initG, queueSize)

		ctx, cancel := context.WithCancel(context.Background())
		stateCh, err := p.State(ctx, time.Millisecond)
		assert.NoError(t, err)

		go func() {
			cancel()
		}()

		for {
			state, ok := <-stateCh
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}

		close(waitCh)
		_, err = p.Shutdown()
		assert.NoError(t, err)
	})

	t.Run("close chan after shutdown", func(t *testing.T) {
		t.Parallel()

		p := runningPool(t, 1, 3)

		ch, err := p.State(context.Background(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err = p.Shutdown()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})

	t.Run("close chan after shutdown now", func(t *testing.T) {
		t.Parallel()

		p := runningPool(t, 1, 3)

		ch, err := p.State(context.Background(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err = p.ShutdownNow()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})
}
