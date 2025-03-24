package plus_sync

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Cond is a condition variable that can be used to wait for a condition to be met.
//
// When change condition, you should use the `Locker` to lock the condition.
//
// In Go memory model, Cond guantees the calling of Broadcast and Signal methods.
//
// In most simple cases, is better to use channels instead of Cond.
// Call broadcast method on a closed channel, Signal should give a channel send message.
type Cond struct {
	notifyList *notifyList
	checker    unsafe.Pointer // pointer to itself to check it's being used by copy.
	once       sync.Once      // use to init notifyList
	Locker     sync.Locker    // Locker is used to Lock when observing or changing condition.
	noCopy     noCopy
}

// NewCond creates a new Cond instance with the provided locker.
func NewCond(locker sync.Locker) *Cond {
	return &Cond{
		Locker: locker,
	}
}

// Wait waits for the condition to be met, unlocking the locker while waiting.
// It locks the locker again after the wait is over.
func (c *Cond) Wait(ctx context.Context) error {
	c.checkCopy()
	c.checkFirstUse()

	node := c.notifyList.add()

	c.Locker.Unlock()
	defer c.Locker.Lock()

	return c.notifyList.wait(ctx, node)
}

// Signal notifies one waiting goroutine that the condition has been met.
func (c *Cond) Signal() {
	c.checkCopy()
	c.checkFirstUse()

	c.notifyList.notifyOne()
}

// Broadcast notifies all waiting goroutines that the condition has been met.
func (c *Cond) Broadcast() {
	c.checkCopy()
	c.checkFirstUse()

	c.notifyList.notifyAll()
}

// checkCopy checks if the Cond instance has been copied and panics if it has.
func (c *Cond) checkCopy() {
	// checking the pointer saved by checker is equal to the current pointer.(not init when create Cond, so it possible to be not equal)
	if c.checker != unsafe.Pointer(c) &&
		// when first time to init, c.checker is zero value.
		// so use CompareAndSwapPointer to set the pointer.
		!atomic.CompareAndSwapPointer(&c.checker, nil, unsafe.Pointer(c)) &&
		// checking the pointer again
		c.checker != unsafe.Pointer(c) {
		panic("[sync] Cond is copied")
	}
}

// checkFirstUse initializes the notifyList if it hasn't been initialized yet.
func (c *Cond) checkFirstUse() {
	c.once.Do(func() {
		c.notifyList = newNotifyList()
	})
}

// notifyList manages a list of waiting notifications using a mutex and a chanList.
type notifyList struct {
	mutex sync.Mutex
	list  *chanList
}

// newNotifyList creates a new notifyList instance.
func newNotifyList() *notifyList {
	return &notifyList{
		mutex: sync.Mutex{},
		list:  newChanList(),
	}
}

// add adds a new channel node to the notifyList.
func (nl *notifyList) add() *chanNode {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	node := nl.list.alloc()
	nl.list.pushBack(node)

	return node
}

// wait waits for a notification or context cancellation.
func (nl *notifyList) wait(ctx context.Context, node *chanNode) error {
	ch := node.Val

	defer nl.list.free(node)

	select {
	case <-ctx.Done():
		nl.mutex.Lock()
		defer nl.mutex.Unlock()

		select {
		// double check.
		case <-ch:
			if nl.list.len() != 0 {
				nl.notifyNext()
			}
		default:
			// if can not receive the signal from the channel, means the channel is never used,
			// waiting object can be removed from the list.
			nl.list.remove(node)
		}
		return ctx.Err()

	case <-ch:
		// receive the signal from the channel means a normal wake-up call
		return nil
	}

}

// notifyNext notifies the next node in the notifyList.
func (nl *notifyList) notifyNext() {
	front := nl.list.front()

	ch := front.Val
	nl.list.remove(front)

	ch <- struct{}{}
}

// notifyOne notifies one waiting node in the notifyList.
func (nl *notifyList) notifyOne() {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	if nl.list.len() == 0 {
		return
	}

	nl.notifyNext()
}

// notifyAll notifies all waiting nodes in the notifyList.
func (nl *notifyList) notifyAll() {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	for nl.list.len() != 0 {
		nl.notifyNext()
	}
}

// chanNode represents a node in the chanList, containing a channel and pointers to previous and next nodes.
type chanNode struct {
	prev *chanNode
	next *chanNode
	Val  chan struct{}
}

// chanList is a doubly linked list used to store channels.
type chanList struct {
	sentinel *chanNode
	size     int
	pool     *sync.Pool
}

// newChanList creates a new chanList instance.
func newChanList() *chanList {
	sentinel := &chanNode{}
	sentinel.prev = sentinel
	sentinel.next = sentinel

	return &chanList{
		sentinel: sentinel,
		size:     0,
		pool: &sync.Pool{
			New: func() any {
				return &chanNode{
					Val: make(chan struct{}, 1),
				}
			},
		},
	}
}

// len returns the number of channels in the chanList.
func (cl *chanList) len() int {
	return cl.size
}

// front returns the first channel node in the chanList.
func (cl *chanList) front() *chanNode {
	return cl.sentinel.next
}

// alloc allocates a new channel node from the pool.
func (cl *chanList) alloc() *chanNode {
	return cl.pool.Get().(*chanNode)
}

// pushBack adds a new channel node to the end of the chanList.
func (cl *chanList) pushBack(node *chanNode) {
	node.next = cl.sentinel
	node.prev = cl.sentinel.prev

	cl.sentinel.prev.next = node
	cl.sentinel.prev = node
	cl.size++
}

// remove removes a channel node from the chanList.
func (cl *chanList) remove(node *chanNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = nil
	node.next = nil
	cl.size--
}

// free returns a channel node to the pool for reuse.
func (cl *chanList) free(node *chanNode) {
	cl.pool.Put(node)
}

// noCopy is a struct used to prevent copying of Cond instances.
type noCopy struct{}

func (*noCopy) Lock() {}

func (*noCopy) Unlock() {}
