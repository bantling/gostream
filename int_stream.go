package stream

import (
	"sort"
)

// intSliceIterator is an iterator for an array
type intSliceIterator struct {
	array []int
	index int
}

// next iterates the array
func (iter *intSliceIterator) next() (int, bool) {
	if iter.index < len(iter.array) {
		next := iter.array[iter.index]
		iter.index++
		return next, true
	}

	return 0, false
}

// IntStream is the int specialization of Stream
type IntStream struct {
	iterator func() (int, bool)
}

// Construct a new IntStream of an iterator
func NewIntStream(iter func() (int, bool)) IntStream {
	return IntStream{iterator: iter}
}

// Construct a new IntStream of an array of values
func NewIntStreamOf(array ...int) IntStream {
	arrayIter := intSliceIterator{array: array}
	return IntStream{iterator: arrayIter.next}
}

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (s IntStream) AllMatch(f func(element int) bool) bool {
	allMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if allMatch = f(next); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s IntStream) AnyMatch(f func(element int) bool) bool {
	anyMatch := false

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if anyMatch = f(next); anyMatch {
			break
		}
	}

	return anyMatch
}

// Average returns an optional average value
func (s IntStream) Average() (float64, bool) {
	var (
		sum int
		count int
	)

	s.ForEach(func(element int) {
		sum += element
		count++
	})

	return float64(sum) / float64(count), count > 0
}

// Concat concatenates two IntStreams into a new IntStream that contains all the elements
// of this IntStream followed by all elements of the IntStream passed
func (s IntStream) Concat(os IntStream) IntStream {
	firstIter := true

	return IntStream{
		iterator: func() (int, bool) {
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
func (s IntStream) Count() int {
	count := 0

	s.ForEach(func(int) { count++ })

	return count
}

// Distinct returns the distinct elements only
func (s IntStream) Distinct() IntStream {
	alreadyRead := map[int]bool{}

	return s.Filter(func(element int) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return true
		}

		return false
	})
}

// Duplicates returns the duplicate elements only
func (s IntStream) Duplicate() IntStream {
	alreadyRead := map[int]bool{}

	return s.Filter(func(element int) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return false
		}

		return true
	})
}

// Filter returns a new IntStream of all elements that pass the given predicate
func (s IntStream) Filter(f func(element int) bool) IntStream {
	return IntStream{
		iterator: func() (int, bool) {
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
func (s IntStream) First() (int, bool) {
	return s.iterator()
}

// ForEach invokes a consumer with each element of the IntStream
func (s IntStream) ForEach(f func(element int)) {
	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		f(next)
	}
}

// GroupBy groups elements by executing the given function on each value to get a key,
// and appending the element to the end of a slice associated with the key in the resulting map.
func (s IntStream) GroupBy(f func(element int) (key interface{})) map[interface{}][]int {
	m := map[interface{}][]int{}

	s.Reduce(
		m,
		func(accumulator interface{}, element int) interface{} {
			k := f(element)
			m[k] = append(m[k], element)
			return m
		},
	)

	return m
}

// Iterate returns a IntStream of an iterative calculation, f(seed), f(f(seed)), ...
func (s IntStream) Iterate(seed int, f func(int) int) IntStream {
	acculumator := seed

	return IntStream{
		iterator: func() (int, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		},
	}
}

// Last returns the optional last element
func (s IntStream) Last() (int, bool) {
	var (
		next    int
		hasNext bool
	)

	s.ForEach(func(element int) {
		next = element
		hasNext = true
	})

	return next, hasNext
}

// Limit returns a new IntStream that only iterates the first n elements, ignoring the rest
func (s IntStream) Limit(n int) IntStream {
	elementsRead := 0
	done := false

	return IntStream{
		iterator: func() (int, bool) {
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

// Map each element to another element
func (s IntStream) Map(f func(element int) int) IntStream {
	return IntStream{
		iterator: func() (int, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// Map each element to a float
func (s IntStream) MapToFloat(f func(element int) float64) FloatStream {
	return FloatStream{
		iterator: func() (float64, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// Map each element to an object
func (s IntStream) MapToObject(f func(element int) interface{}) Stream {
	return Stream{
		iterator: func() (interface{}, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return nil, false
		},
	}
}

// Map each element to a string
func (s IntStream) MapToString(f func(element int) string) StringStream {
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
func (s IntStream) Max() (int, bool) {
	max, hasMax := s.iterator()
	if hasMax {
		s.ForEach(func(element int) {
			if max < element {
				max = element
			}
		})
	}

	return max, hasMax
}

// Min returns an optional minimum value according to the provided comparator
func (s IntStream) Min() (int, bool) {
	min, hasMin := s.iterator()
	if hasMin {
		s.ForEach(func(element int) {
			if element < min {
				min = element
			}
		})
	}

	return min, hasMin
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (s IntStream) NoneMatch(f func(element int) bool) bool {
	noneMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if noneMatch = !f(next); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Peek calls a function to examine each value and perform an additional operation
func (s IntStream) Peek(f func(int)) IntStream {
	return IntStream{
		iterator: func() (int, bool) {
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
func (s IntStream) Reduce(
	identity interface{},
	f func(accumulator interface{}, element int) interface{},
) interface{} {
	result := identity

	s.ForEach(func(element int) {
		result = f(result, element)
	})

	return result
}

// Skip returns a new IntStream that skips the first n elements
func (s IntStream) Skip(n int) IntStream {
	done := false

	return IntStream{
		iterator: func() (int, bool) {
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

// Sorted returns a new IntStream with the values sorted by the provided comparator..
func (s IntStream) Sorted() IntStream {
	var sortedIter func() (int, bool)
	done := false

	return IntStream{
		iterator: func() (int, bool) {
			if !done {
				// Sort all IntStream elements
				sorted := s.ToSlice()
				sort.Ints(sorted)

				sortedIter = (&intSliceIterator{array: sorted}).next
				done = true
			}

			// Return next sorted element
			return sortedIter()
		},
	}
}

// Sum returns an optional sum value
func (s IntStream) Sum() (int, bool) {
	var (
		sum int
		haveSum bool
	)

	s.ForEach(func(element int) {
		sum += element
		haveSum = true
	})

	return sum, haveSum
}

// ToMap returns a map of all elements by invoking the given function to a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s IntStream) ToMap(f func(int) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	s.ForEach(func(element int) {
		k, v := f(element)
		m[k] = v
	})

	return m
}

// ToSlice returns a slice of all elements
func (s IntStream) ToSlice() []int {
	var array []int

	s.ForEach(func(element int) {
		array = append(array, element)
	})

	return array
}
