package stream

import (
	"sort"
)

// arrayIterator is an iterator for an array
type arrayIterator struct {
	array []interface{}
	index int
}

// Next iterates the array
func (iter *arrayIterator) next() (interface{}, bool) {
	if iter.index < len(iter.array) {
		next := iter.array[iter.index]
		iter.index++
		return next, true
	}

	return nil, false
}

// Stream is the base object type for streams.
// It is based on a simple iterator function that returns (interface{}, bool).
// If the bool result is true, the interface{} is a valid value, and there may be more values.
// If the bool result is false, the interface{} is not a valid value, and there are no more values.
// An iterator can be infinite by simply never returning false.
//
// Some methods are chaining methods, they return a new stream
// Some functions are terminal, they return a non-stream result
// Some terminal functions return optional values by returning (<type>, bool),
// the meaning of the results is the same as for an iterator.
// Some functions accept a comparator.
// Normally, a comparator takes two int slice indexes and returns true if the value at index1 < the value at index2.
// For a Stream, a comparator takes two elements and returns true if element1 < element 2.
type Stream struct {
	iterator func() (interface{}, bool)
}

// Construct a new stream of an iterator
func NewStream(iter func() (interface{}, bool)) Stream {
	return Stream{iterator: iter}
}

// Construct a new stream of an array of values
func NewStreamOf(array ...interface{}) Stream {
	arrayIter := arrayIterator{array: array}
	return Stream{iterator: arrayIter.next}
}

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (s Stream) AllMatch(f func(element interface{}) bool) bool {
	allMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if allMatch = f(next); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s Stream) AnyMatch(f func(element interface{}) bool) bool {
	anyMatch := false

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if anyMatch = f(next); anyMatch {
			break
		}
	}

	return anyMatch
}

// Concat concatenates two streams into a new stream that contains all the elements
// of this stream followed by all elements of the stream passed
func (s Stream) Concat(os Stream) Stream {
	firstIter := true

	return Stream{
		iterator: func() (interface{}, bool) {
			if firstIter {
				if next, hasNext := s.iterator(); hasNext {
					return next, hasNext
				}

				// Switch to second iterator and return first element
				firstIter = false
				return os.iterator()
			} else {
				return os.iterator()
			}
		},
	}
}

// Count returns the count of all elements
func (s Stream) Count() int {
	count := 0

	s.ForEach(func(interface{}) { count++ })

	return count
}

// Distinct returns the distinct elements only
func (s Stream) Distinct() Stream {
	alreadyRead := map[interface{}]bool{}

	return s.Filter(func(val interface{}) bool {
		if !alreadyRead[val] {
			alreadyRead[val] = true
			return true
		}

		return false
	})
}

// Duplicates returns the duplicate elements only
func (s Stream) Duplicate() Stream {
	alreadyRead := map[interface{}]bool{}

	return s.Filter(func(val interface{}) bool {
		if !alreadyRead[val] {
			alreadyRead[val] = true
			return false
		}

		return true
	})
}

// Filter returns a new stream of all elements that pass the given predicate
func (s Stream) Filter(f func(element interface{}) bool) Stream {
	return Stream{
		iterator: func() (interface{}, bool) {
			for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
				if f(next) {
					return next, true
				}
			}

			return nil, false
		},
	}
}

// First returns the optional first element
func (s Stream) First() (interface{}, bool) {
	return s.iterator()
}

// ForEach invokes a consumer with each element of the stream
func (s Stream) ForEach(f func(element interface{})) {
	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		f(next)
	}
}

// GroupBy groups elements by executing the given function on each value to get a key,
// and appending the element to the end of a slice associated with the key in the resulting map.
func (s Stream) GroupBy(f func(element interface{}) (key interface{})) map[interface{}][]interface{} {
	m := map[interface{}][]interface{}{}

	s.ForEach(func(val interface{}) {
		k := f(val)
		m[k] = append(m[k], val)
	})

	return m
}

// Iterate returns a stream of an iterative calculation, f(seed), f(f(seed)), ...
func (s Stream) Iterate(seed interface{}, f func(interface{}) interface{}) Stream {
	first := true
	var acculumator interface{}

	return Stream{
		iterator: func() (interface{}, bool) {
			if first {
				first = false
				acculumator = f(seed)
			} else {
				acculumator = f(acculumator)
			}

			return acculumator, true
		},
	}
}

// Last returns the optional last element
func (s Stream) Last() (interface{}, bool) {
	var (
		next    interface{}
		hasNext bool
	)

	s.ForEach(func(val interface{}) {
		next = val
		hasNext = true
	})

	return next, hasNext
}

// Limit returns a new stream that only iterates the first n elements, ignoring the rest
func (s Stream) Limit(n int) Stream {
	elementsRead := 0
	done := false

	return Stream{
		iterator: func() (interface{}, bool) {
			if done {
				return nil, false
			}

			next, hasNext := s.iterator()
			if !hasNext {
				done = true
				return nil, false
			}

			elementsRead++
			done = elementsRead == n
			return next, hasNext
		},
	}
}

// Map each element to a new element, possibly of a different type
func (s Stream) Map(f func(element interface{}) interface{}) Stream {
	return Stream{
		iterator: func() (interface{}, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return nil, false
		},
	}
}

// Max returns an optional maximum value according to the provided comparator
func (s Stream) Max(less func(element1, element2 interface{}) bool) (interface{}, bool) {
	max, hasMax := s.iterator()
	if hasMax {
		s.ForEach(func(val interface{}) {
			if less(max, val) {
				max = val
			}
		})
	}

	return max, hasMax
}

// Min returns an optional minimum value according to the provided comparator
func (s Stream) Min(less func(element1, element2 interface{}) bool) (interface{}, bool) {
	min, hasMin := s.iterator()
	if hasMin {
		s.ForEach(func(val interface{}) {
			if less(val, min) {
				min = val
			}
		})
	}

	return min, hasMin
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (s Stream) NoneMatch(f func(element interface{}) bool) bool {
	noneMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if noneMatch = !f(next); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Peek calls a function to examine each value and perform an additional operation
func (s Stream) Peek(f func(interface{})) Stream {
	return Stream{
		iterator: func() (interface{}, bool) {
			next, hasNext := s.First()
			if hasNext {
				f(next)
			}

			return next, hasNext
		},
	}
}

// Reduce uses a function to reduce the stream to a single optional value,
// depending on the number of elements, and whether or not the optional identity is provided.
//
// If no identity provided:
// 0 elements = invalid
// 1 element = element1
// 2 elements = f(element1, element2)
// 3 or more elements = f(f(element1, element2), val3)...
//
// If an identity is provided:
// 0 elements = identity
// 1 element = f(identity, element1)
// 2 or more elements = f(f(identity, element1), element2)...
func (s Stream) Reduce(
	f func(element1, element2 interface{}) interface{},
	identityOpt ...interface{},
) (interface{}, bool) {
	var (
		result interface{}
		valid  bool
	)

	if len(identityOpt) == 0 {
		result, valid = s.iterator()
	} else {
		result = identityOpt[0]
		valid = true
	}

	s.ForEach(func(next interface{}) {
		result = f(result, next)
	})

	return result, valid
}

// Skip returns a new stream that skips the first n elements
func (s Stream) Skip(n int) Stream {
	done := false

	return Stream{
		iterator: func() (interface{}, bool) {
			// Skip n elements only once
			if !done {
				for i := 1; i <= n; i++ {
					if _, hasNext := s.iterator(); !hasNext {
						break
					}
				}

				done = true
			}

			// Return next element
			return s.iterator()
		},
	}
}

// Sorted returns a new stream with the values sorted by the provided comparator..
func (s Stream) Sorted(less func(element1, element2 interface{}) bool) Stream {
	var sortedIter func() (interface{}, bool)
	done := false

	return Stream{
		iterator: func() (interface{}, bool) {
			if !done {
				// Sort all stream elements
				sorted := s.ToSlice()
				sort.Slice(sorted, func(i, j int) bool {
					return less(sorted[i], sorted[j])
				})

				sortedIter = (&arrayIterator{array: sorted}).next
				done = true
			}

			// Return first sorted element
			return sortedIter()
		},
	}
}

// ToMap returns a map of all elements by invoking the given function to a key/value pair for the map
func (s Stream) ToMap(f func(interface{}) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	s.ForEach(func(val interface{}) {
		k, v := f(val)
		m[k] = v
	})

	return m
}

// ToSlice returns a slice of all elements
func (s Stream) ToSlice() []interface{} {
	var array []interface{}

	s.ForEach(func(val interface{}) {
		array = append(array, val)
	})

	return array
}
