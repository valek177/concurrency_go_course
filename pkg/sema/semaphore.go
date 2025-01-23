package sema

import (
	"sync"
)

type Semaphore struct {
	cond  *sync.Cond
	count int
	max   int
}

func NewSemaphore(max int) *Semaphore {
	mutex := &sync.Mutex{}
	return &Semaphore{
		max:  max,
		cond: sync.NewCond(mutex),
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.count >= s.max {
		s.cond.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--
	s.cond.Signal()
}
