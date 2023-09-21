// lightweight implementation of circuit breaker pattern
package fuse

import (
	"fmt"
	"sync"
	"time"
)

// sender interface exposes a send method
//
// Breaker uses this pattern to trigger send call - see test example
type Sender interface {
	Send() bool
}

// BreakerState type - int
type BreakerState int

// this breaker struct implements circuit-breaker pattern
type Breaker struct {
	debug           bool
	mu              sync.Mutex
	failureCount    int
	state           BreakerState
	maxFailureCount int
	resetTimeout    int
}

// sets debug flag to true for more information in console
func (b *Breaker) Debug() *Breaker {
	b.debug = true
	fmt.Println("breaker.maxfailurecount", b.getFailureCount())
	fmt.Println("breaker.resetcount", b.getResetTimeout())
	return b
}

// sets breaker max failure count
func (b *Breaker) SetMaxFailureCount(value int) {
	b.maxFailureCount = value
}

// sets breaker reset timeout
func (b *Breaker) SetResetTimeout(value int) {
	b.resetTimeout = value
}

// try send request through breaker
//
// returns true if request was successful
func (b *Breaker) Send(sender Sender) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == OPEN {
		if b.debug {
			fmt.Println("circuit is broken/open")
		}
		// Circuit is open, do not allow requests.
		return false
	}

	out := sender.Send()
	if out {
		b.failureCount = 0
		return true
	}

	// Simulate a failed request.
	b.failureCount++
	if b.failureCount >= b.getFailureCount() {
		b.setState(OPEN)
		go b.Reset()
	}

	return false
}

// update current breaker state with given state
func (b *Breaker) setState(bs BreakerState) {
	if b.debug {
		fmt.Println("current state", b.state)
		if b.state == bs {
			fmt.Println("no change in state")
		} else {
			fmt.Println("next state", bs)
		}
	}
	b.state = bs
}

// resets the breaker circuit after a timeout.
//
// returns breaker channel
func (b *Breaker) reset(h chan<- *Breaker) {
	defer close(h)
	time.Sleep(time.Duration(b.getResetTimeout()) * time.Second)
	b.mu.Lock()
	if b.debug {
		fmt.Println("reset triggered")
	}
	b.setState(CLOSED)
	b.mu.Unlock()
	h <- b
}

// resets the breaker circuit after a timeout.
//
// returns breaker channel
func (b *Breaker) Reset() <-chan *Breaker {
	br := make(chan *Breaker)
	go b.reset(br)
	return br
}

// returns default or specified failure count
func (b *Breaker) getFailureCount() int {
	if b.maxFailureCount > 0 {
		return b.maxFailureCount
	} else {
		return defaultMaxFailureCount
	}
}

// returns default or specified reset timeout
func (b *Breaker) getResetTimeout() int {
	if b.resetTimeout > 0 {
		return b.resetTimeout
	} else {
		return defaultMaxTimeout
	}
}
