package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func poolInternalState(p *BlockTaskPool) int32 {
	for {
		state := atomic.LoadInt32(&p.state)
		if state != stateLocked {
			return state
		}
	}
}

func TestTaskPool_NewBlockTaskPool(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		initG     int32
		queueSize int32
		wantErr   error
	}{
		{
			name:      "basic",
			initG:     1,
			queueSize: 1,
			wantErr:   nil,
		}, {
			name:      "queue size is negative",
			initG:     1,
			queueSize: -1,
			wantErr:   errInvalidParam,
		}, {
			name:      "queue size is 0",
			initG:     1,
			queueSize: 0,
			wantErr:   nil,
		}, {
			name:      "queue size greater than 0",
			initG:     1,
			queueSize: 1,
			wantErr:   nil,
		}, {
			name:      "init goroutines is negative",
			initG:     -1,
			queueSize: 1,
			wantErr:   errInvalidParam,
		}, {
			name:      "init goroutines is 0",
			initG:     0,
			queueSize: 1,
			wantErr:   errInvalidParam,
		}, {
			name:      "init goroutines greater than 0",
			initG:     1,
			queueSize: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p, err := NewBlockTaskPool(tc.initG, tc.queueSize)
			assert.ErrorIs(t, err, tc.wantErr)

			if err == nil {
				assert.NotNil(t, p)
				assert.Equal(t, stateCreated, poolInternalState(p))

				assert.Equal(t, tc.initG, p.initG)
				assert.Equal(t, int(tc.queueSize), cap(p.queue))
			}
		})
	}

	t.Run("with max idle time", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithMaxIdleTime(time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(1), p.coreG)
		assert.Equal(t, int32(1), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, time.Second, p.maxIdleTime)
	})

	t.Run("with submit timeout", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithSubmitTimeout(time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(1), p.coreG)
		assert.Equal(t, int32(1), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, time.Second, p.submitTimeout)
	})

	t.Run("with core goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with max goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithMaxG(2))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with core and max goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2), WithMaxG(4))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(4), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, int32(4), p.maxG)
	})

	t.Run("with core != init and max == init", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2), WithMaxG(1))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with core == init and max != init", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(1), WithMaxG(4))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(4), p.coreG)
		assert.Equal(t, int32(4), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with queue backlog rate", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(0.5))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, 0.5, p.queueBacklogRate)
	})

	t.Run("queue backlog rate is negative", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(-1))
		assert.ErrorIs(t, err, errInvalidParam)
		assert.Nil(t, p)
	})

	t.Run("queue backlog rate is greater than 1", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(1.1))
		assert.ErrorIs(t, err, errInvalidParam)
		assert.Nil(t, p)
	})
}

func TestTaskPool_Submit(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name       string
		poolFunc   func(t *testing.T) *BlockTaskPool
		submitFunc func(t *testing.T, p *BlockTaskPool)
		wantErr    error
	}{
		{
			name: "basic",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				p, err := NewBlockTaskPool(1, 3)
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				var err error
				for range 3 {
					err = p.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						return nil
					}))
					assert.NoError(t, err)
				}
			},
		}, {
			name: "nil task",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				p, err := NewBlockTaskPool(1, 1)
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				err := p.Submit(context.Background(), nil)
				assert.ErrorIs(t, err, errInvalidTask)
			},
		}, {
			name: "submit timeout",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				p, err := NewBlockTaskPool(1, 1, WithSubmitTimeout(time.Millisecond))
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				done := make(chan struct{})

				err := p.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-done
					return nil
				}))
				assert.NoError(t, err)

				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
				err = p.Submit(ctx, TaskFunc(func(ctx context.Context) error {
					<-done
					return nil
				}))
				cancel()
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				close(done)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := tc.poolFunc(t)
			assert.NotNil(t, p)
			assert.Equal(t, poolInternalState(p), stateCreated)

			tc.submitFunc(t, p)
		})
	}
}

func TestTaskPool_Shutdown(t *testing.T) {
	t.Parallel()

	p1, err := NewBlockTaskPool(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p1)
	assert.Equal(t, poolInternalState(p1), stateCreated)

	done, err := p1.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, errPoolIsNotRunning)
	assert.Equal(t, poolInternalState(p1), stateCreated)

	err = p1.Start()
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p1), stateRunning)

	done, err = p1.Shutdown()
	assert.NotNil(t, done)
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p1), stateClosing)
	<-done
	assert.Equal(t, poolInternalState(p1), stateClosed)

	done, err = p1.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, errPoolIsClosed)

	p2, err := NewBlockTaskPool(1, 3)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
	assert.Equal(t, poolInternalState(p2), stateCreated)

	tasks, err := p2.ShutdownNow()
	assert.Equal(t, 0, len(tasks))
	assert.ErrorIs(t, err, errPoolIsNotRunning)
	assert.Equal(t, poolInternalState(p2), stateCreated)

	err = p2.Start()
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p2), stateRunning)

	tasks, err = p2.ShutdownNow()
	assert.Equal(t, 0, len(tasks))
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p2), stateClosed)

	_, err = p2.ShutdownNow()
	assert.ErrorIs(t, err, errPoolIsClosed)
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
