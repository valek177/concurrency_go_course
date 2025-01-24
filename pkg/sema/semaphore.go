package sema

import (
	"sync"
)

// Semaphore is a struct for semaphore
type Semaphore struct {
	cond  *sync.Cond
	count int
	max   int
}

// NewSemaphore returns new semaphore
func NewSemaphore(max int) *Semaphore {
	mutex := &sync.Mutex{}
	return &Semaphore{
		max:  max,
		cond: sync.NewCond(mutex),
	}
}

// Acquire acquires semaphore
func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.count >= s.max {
		s.cond.Wait()
	}

	s.count++
}

// Release releases semaphore
func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--
	s.cond.Signal()
}
