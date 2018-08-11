package set

// Add a variable number of elements to the set.
func (s *Set) Add(elements ...interface{}) {
	if len(elements) == 0 {
		return
	}

	s.mutex.Lock()
	for _, element := range elements {
		s.m[element] = empty
	}
	s.mutex.Unlock()
}

// Remove a variable number of elements from the set.
func (s *Set) Remove(elements ...interface{}) {
	if len(elements) == 0 {
		return
	}

	s.mutex.Lock()
	for _, element := range elements {
		delete(s.m, element)
	}
	s.mutex.Unlock()
}

// Has returns true if all given variable number of elements is included in
// the set. If not all elements are included or there were no elements given,
// false is returned.
func (s *Set) Has(elements ...interface{}) (has bool) {
	if len(elements) == 0 {
		return false
	}

	s.mutex.RLock()
	for _, element := range elements {
		if _, has = s.m[element]; !has {
			break
		}
	}
	s.mutex.RUnlock()

	return
}

// Size returns the number of elements in a set.
func (s *Set) Size() (l int) {
	s.mutex.RLock()
	l = len(s.m)
	s.mutex.RUnlock()

	return
}

// Clear empties out a set.
func (s *Set) Clear() {
	s.mutex.Lock()
	s.m = make(map[interface{}]struct{})
	s.mutex.Unlock()
}

// StringSlice returns a slice of strings of the set items.
func (s *Set) StringSlice() []string {
	slice := make([]string, 0, s.Size())

	s.mutex.Lock()
	for k := range s.m {
		slice = append(slice, k.(string))
	}
	s.mutex.Unlock()

	return slice
}
