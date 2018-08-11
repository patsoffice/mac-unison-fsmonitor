package set

import "sync"

// Interface defines the supported thread operations
type Interface interface {
	Add(elements ...interface{})
	Remove(elements ...interface{})
	Has(elements ...interface{}) bool
	Size() int
	Clear()
	Each(func(item interface{}) bool)
}

var empty = struct{}{}

// Set defines a thread safe set data structure that has very limited
// functionality.
type Set struct {
	m     map[interface{}]struct{} // struct{} doesn't take up space
	mutex *sync.RWMutex
}
