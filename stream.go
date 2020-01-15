package stream

import (
	"sort"
)

// sliceIterator is an iterator for an array
type sliceIterator struct {
	array []interface{}
	index int
}

// next iterates the array
func (iter *sliceIterator) next() (interface{}, bool) {
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
	arrayIter := sliceIterator{array: array}
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

	return s.Filter(func(element interface{}) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return true
		}

		return false
	})
}

// Duplicates returns the duplicate elements only
func (s Stream) Duplicate() Stream {
	alreadyRead := map[interface{}]bool{}

	return s.Filter(func(element interface{}) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
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

	s.Reduce(
		m,
		func(accumulator interface{}, element interface{}) interface{} {
			k := f(element)
			m[k] = append(m[k], element)
			return m
		},
	)

	return m
}

// Iterate returns a stream of an iterative calculation, f(seed), f(f(seed)), ...
func (s Stream) Iterate(seed interface{}, f func(interface{}) interface{}) Stream {
	acculumator := seed

	return Stream{
		iterator: func() (interface{}, bool) {
			acculumator = f(acculumator)

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

	s.ForEach(func(element interface{}) {
		next = element
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

// MapToInt each element to an int
func (s Stream) MapToInt(f func(element interface{}) int) IntStream {
	return IntStream{
		iterator: func() (int, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// MapToFloat each element to a float
func (s Stream) MapToFloat(f func(element interface{}) float64) FloatStream {
	return FloatStream{
		iterator: func() (float64, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// MapToInt each element to an int
func (s Stream) MapToString(f func(element interface{}) string) StringStream {
	return StringStream{
		iterator: func() (string, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return "", false
		},
	}
}

// Max returns an optional maximum value according to the provided comparator
func (s Stream) Max(less func(element1, element2 interface{}) bool) (interface{}, bool) {
	max, hasMax := s.iterator()
	if hasMax {
		s.ForEach(func(element interface{}) {
			if less(max, element) {
				max = element
			}
		})
	}

	return max, hasMax
}

// Min returns an optional minimum value according to the provided comparator
func (s Stream) Min(less func(element1, element2 interface{}) bool) (interface{}, bool) {
	min, hasMin := s.iterator()
	if hasMin {
		s.ForEach(func(element interface{}) {
			if less(element, min) {
				min = element
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
			next, hasNext := s.iterator()
			if hasNext {
				f(next)
			}

			return next, hasNext
		},
	}
}

// Reduce uses a function to reduce the stream to a single value by iteratively executing a function
// with the current accumulated value and the next stream element.
// The identity provided is the initial accumulated value, which means the result type is the
// same type as the initial value, which can be any type.
// If there are no elements in the strea, the result is the identity.
// Otherwise, the result is f(f(identity, element1), element2)...
func (s Stream) Reduce(
	identity interface{},
	f func(accumulator interface{}, element2 interface{}) interface{},
) interface{} {
	result := identity

	s.ForEach(func(element interface{}) {
		result = f(result, element)
	})

	return result
}

// ReverseSorted returns a stream with elements sorted in decreasing order.
// The provided function must compare elements in increasing order, same as for Sorted.
func (s Stream) ReverseSorted(less func(element1, element2 interface{}) bool) Stream {
	return s.Sorted(func(element1, element2 interface{}) bool {
		return !less(element1, element2)
	})
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

				sortedIter = (&sliceIterator{array: sorted}).next
				done = true
			}

			// Return next sorted element
			return sortedIter()
		},
	}
}

// ToMap returns a map of all elements by invoking the given function to get a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s Stream) ToMap(f func(interface{}) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	s.ForEach(func(element interface{}) {
		k, v := f(element)
		m[k] = v
	})

	return m
}

// ToSlice returns a slice of all elements
func (s Stream) ToSlice() []interface{} {
	var array []interface{}

	s.ForEach(func(element interface{}) {
		array = append(array, element)
	})

	return array
}
