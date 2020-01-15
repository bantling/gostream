package stream

import (
	"sort"
)

// floatSliceIterator is an iterator for an array
type floatSliceIterator struct {
	array []float64
	index int
}

// next iterates the array
func (iter *floatSliceIterator) next() (float64, bool) {
	if iter.index < len(iter.array) {
		next := iter.array[iter.index]
		iter.index++
		return next, true
	}

	return 0, false
}

// FloatStream is the float64 specialization of Stream
type FloatStream struct {
	iterator func() (float64, bool)
}

// Construct a new FloatStream of an iterator
func NewFloatStream(iter func() (float64, bool)) FloatStream {
	return FloatStream{iterator: iter}
}

// Construct a new FloatStream of an array of values
func NewFloatStreamOf(array ...float64) FloatStream {
	arrayIter := floatSliceIterator{array: array}
	return FloatStream{iterator: arrayIter.next}
}

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (s FloatStream) AllMatch(f func(element float64) bool) bool {
	allMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if allMatch = f(next); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s FloatStream) AnyMatch(f func(element float64) bool) bool {
	anyMatch := false

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if anyMatch = f(next); anyMatch {
			break
		}
	}

	return anyMatch
}

// Average returns an optional average value
func (s FloatStream) Average() (float64, bool) {
	var (
		sum   float64
		count int
	)

	s.ForEach(func(element float64) {
		sum += element
		count++
	})

	return float64(sum) / float64(count), count > 0
}

// Concat concatenates two FloatStreams into a new FloatStream that contains all the elements
// of this FloatStream followed by all elements of the FloatStream passed
func (s FloatStream) Concat(os FloatStream) FloatStream {
	firstIter := true

	return FloatStream{
		iterator: func() (float64, bool) {
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
func (s FloatStream) Count() int {
	count := 0

	s.ForEach(func(float64) { count++ })

	return count
}

// Distinct returns the distinct elements only
func (s FloatStream) Distinct() FloatStream {
	alreadyRead := map[float64]bool{}

	return s.Filter(func(element float64) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return true
		}

		return false
	})
}

// Duplicates returns the duplicate elements only
func (s FloatStream) Duplicate() FloatStream {
	alreadyRead := map[float64]bool{}

	return s.Filter(func(element float64) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return false
		}

		return true
	})
}

// Filter returns a new FloatStream of all elements that pass the given predicate
func (s FloatStream) Filter(f func(element float64) bool) FloatStream {
	return FloatStream{
		iterator: func() (float64, bool) {
			for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
				if f(next) {
					return next, true
				}
			}

			return 0, false
		},
	}
}

// First returns the optional first element
func (s FloatStream) First() (float64, bool) {
	return s.iterator()
}

// ForEach invokes a consumer with each element of the FloatStream
func (s FloatStream) ForEach(f func(element float64)) {
	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		f(next)
	}
}

// GroupBy groups elements by executing the given function on each value to get a key,
// and appending the element to the end of a slice associated with the key in the resulting map.
func (s FloatStream) GroupBy(f func(element float64) (key interface{})) map[interface{}][]float64 {
	m := map[interface{}][]float64{}

	s.Reduce(
		m,
		func(accumulator interface{}, element float64) interface{} {
			k := f(element)
			m[k] = append(m[k], element)
			return m
		},
	)

	return m
}

// Iterate returns a FloatStream of an iterative calculation, f(seed), f(f(seed)), ...
func (s FloatStream) Iterate(seed float64, f func(float64) float64) FloatStream {
	acculumator := seed

	return FloatStream{
		iterator: func() (float64, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		},
	}
}

// Last returns the optional last element
func (s FloatStream) Last() (float64, bool) {
	var (
		next    float64
		hasNext bool
	)

	s.ForEach(func(element float64) {
		next = element
		hasNext = true
	})

	return next, hasNext
}

// Limit returns a new FloatStream that only iterates the first n elements, ignoring the rest
func (s FloatStream) Limit(n int) FloatStream {
	elementsRead := 0
	done := false

	return FloatStream{
		iterator: func() (float64, bool) {
			if done {
				return 0, false
			}

			next, hasNext := s.iterator()
			if !hasNext {
				done = true
				return 0, false
			}

			elementsRead++
			done = elementsRead == n
			return next, hasNext
		},
	}
}

// Map each element to a new element, possibly of a different type
func (s FloatStream) Map(f func(element float64) interface{}) Stream {
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
func (s FloatStream) Max() (float64, bool) {
	max, hasMax := s.iterator()
	if hasMax {
		s.ForEach(func(element float64) {
			if max < element {
				max = element
			}
		})
	}

	return max, hasMax
}

// Min returns an optional minimum value according to the provided comparator
func (s FloatStream) Min() (float64, bool) {
	min, hasMin := s.iterator()
	if hasMin {
		s.ForEach(func(element float64) {
			if element < min {
				min = element
			}
		})
	}

	return min, hasMin
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (s FloatStream) NoneMatch(f func(element float64) bool) bool {
	noneMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if noneMatch = !f(next); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Peek calls a function to examine each value and perform an additional operation
func (s FloatStream) Peek(f func(float64)) FloatStream {
	return FloatStream{
		iterator: func() (float64, bool) {
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
func (s FloatStream) Reduce(
	identity interface{},
	f func(accumulator interface{}, element float64) interface{},
) interface{} {
	result := identity

	s.ForEach(func(element float64) {
		result = f(result, element)
	})

	return result
}

// Skip returns a new FloatStream that skips the first n elements
func (s FloatStream) Skip(n int) FloatStream {
	done := false

	return FloatStream{
		iterator: func() (float64, bool) {
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

// Sorted returns a new FloatStream with the values sorted by the provided comparator..
func (s FloatStream) Sorted() FloatStream {
	var sortedIter func() (float64, bool)
	done := false

	return FloatStream{
		iterator: func() (float64, bool) {
			if !done {
				// Sort all FloatStream elements
				sorted := s.ToSlice()
				sort.Float64s(sorted)

				sortedIter = (&floatSliceIterator{array: sorted}).next
				done = true
			}

			// Return next sorted element
			return sortedIter()
		},
	}
}

// Sum returns an optional sum value
func (s FloatStream) Sum() (float64, bool) {
	var (
		sum     float64
		haveSum bool
	)

	s.ForEach(func(element float64) {
		sum += element
		haveSum = true
	})

	return sum, haveSum
}

// ToMap returns a map of all elements by invoking the given function to a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s FloatStream) ToMap(f func(float64) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	s.ForEach(func(element float64) {
		k, v := f(element)
		m[k] = v
	})

	return m
}

// ToSlice returns a slice of all elements
func (s FloatStream) ToSlice() []float64 {
	var array []float64

	s.ForEach(func(element float64) {
		array = append(array, element)
	})

	return array
}
