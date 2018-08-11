package set

import (
	"sync"
)

// New creates and initialize a new set. It accepts a variable number of
// arguments that will make up the initial set of elements.
func New(elements ...interface{}) (s *Set) {
	s = &Set{}
	s.m = make(map[interface{}]struct{})
	s.mutex = &sync.RWMutex{}

	s.Add(elements...)
	return
}
