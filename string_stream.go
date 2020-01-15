package stream

import (
	"sort"
)

// stringSliceIterator is an iterator for an array
type stringSliceIterator struct {
	array []string
	index int
}

// next iterates the array
func (iter *stringSliceIterator) next() (string, bool) {
	if iter.index < len(iter.array) {
		next := iter.array[iter.index]
		iter.index++
		return next, true
	}

	return "", false
}

// StringStream is the string specialization of Stream
type StringStream struct {
	iterator func() (string, bool)
}

// Construct a new StringStream of an iterator
func NewStringStream(iter func() (string, bool)) StringStream {
	return StringStream{iterator: iter}
}

// Construct a new StringStream of an array of values
func NewStringStreamOf(array ...string) StringStream {
	arrayIter := stringSliceIterator{array: array}
	return StringStream{iterator: arrayIter.next}
}

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (s StringStream) AllMatch(f func(element string) bool) bool {
	allMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if allMatch = f(next); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s StringStream) AnyMatch(f func(element string) bool) bool {
	anyMatch := false

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if anyMatch = f(next); anyMatch {
			break
		}
	}

	return anyMatch
}

// Concat concatenates two StringStreams into a new StringStream that contains all the elements
// of this StringStream followed by all elements of the StringStream passed
func (s StringStream) Concat(os StringStream) StringStream {
	firstIter := true

	return StringStream{
		iterator: func() (string, bool) {
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
func (s StringStream) Count() int {
	count := 0

	s.ForEach(func(string) { count++ })

	return count
}

// Distinct returns the distinct elements only
func (s StringStream) Distinct() StringStream {
	alreadyRead := map[string]bool{}

	return s.Filter(func(element string) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return true
		}

		return false
	})
}

// Duplicates returns the duplicate elements only
func (s StringStream) Duplicate() StringStream {
	alreadyRead := map[string]bool{}

	return s.Filter(func(element string) bool {
		if !alreadyRead[element] {
			alreadyRead[element] = true
			return false
		}

		return true
	})
}

// Filter returns a new StringStream of all elements that pass the given predicate
func (s StringStream) Filter(f func(element string) bool) StringStream {
	return StringStream{
		iterator: func() (string, bool) {
			for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
				if f(next) {
					return next, true
				}
			}

			return "", false
		},
	}
}

// First returns the optional first element
func (s StringStream) First() (string, bool) {
	return s.iterator()
}

// ForEach invokes a consumer with each element of the StringStream
func (s StringStream) ForEach(f func(element string)) {
	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		f(next)
	}
}

// GroupBy groups elements by executing the given function on each value to get a key,
// and appending the element to the end of a slice associated with the key in the resulting map.
func (s StringStream) GroupBy(f func(element string) (key interface{})) map[interface{}][]string {
	m := map[interface{}][]string{}

	s.Reduce(
		m,
		func(accumulator interface{}, element string) interface{} {
			k := f(element)
			m[k] = append(m[k], element)
			return m
		},
	)

	return m
}

// Iterate returns a StringStream of an iterative calculation, f(seed), f(f(seed)), ...
func (s StringStream) Iterate(seed string, f func(string) string) StringStream {
	acculumator := seed

	return StringStream{
		iterator: func() (string, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		},
	}
}

// Last returns the optional last element
func (s StringStream) Last() (string, bool) {
	var (
		next    string
		hasNext bool
	)

	s.ForEach(func(element string) {
		next = element
		hasNext = true
	})

	return next, hasNext
}

// Limit returns a new StringStream that only iterates the first n elements, ignoring the rest
func (s StringStream) Limit(n int) StringStream {
	elementsRead := 0
	done := false

	return StringStream{
		iterator: func() (string, bool) {
			if done {
				return "", false
			}

			next, hasNext := s.iterator()
			if !hasNext {
				done = true
				return "", false
			}

			elementsRead++
			done = elementsRead == n
			return next, hasNext
		},
	}
}

// Map each element to a new element
func (s StringStream) Map(f func(element string) string) StringStream {
	return StringStream{
		iterator: func() (string, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return "", false
		},
	}
}

// Map each element to a float
func (s StringStream) MapToFloat(f func(element string) float64) FloatStream {
	return FloatStream{
		iterator: func() (float64, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// Map each element to an int
func (s StringStream) MapToInt(f func(element string) int) IntStream {
	return IntStream{
		iterator: func() (int, bool) {
			if next, hasNext := s.iterator(); hasNext {
				return f(next), true
			}

			return 0, false
		},
	}
}

// Map each element to an object
func (s StringStream) MapToObject(f func(element string) interface{}) Stream {
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
func (s StringStream) Max() (string, bool) {
	max, hasMax := s.iterator()
	if hasMax {
		s.ForEach(func(element string) {
			if max < element {
				max = element
			}
		})
	}

	return max, hasMax
}

// Min returns an optional minimum value according to the provided comparator
func (s StringStream) Min() (string, bool) {
	min, hasMin := s.iterator()
	if hasMin {
		s.ForEach(func(element string) {
			if element < min {
				min = element
			}
		})
	}

	return min, hasMin
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (s StringStream) NoneMatch(f func(element string) bool) bool {
	noneMatch := true

	for next, hasNext := s.iterator(); hasNext; next, hasNext = s.iterator() {
		if noneMatch = !f(next); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Peek calls a function to examine each value and perform an additional operation
func (s StringStream) Peek(f func(string)) StringStream {
	return StringStream{
		iterator: func() (string, bool) {
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
func (s StringStream) Reduce(
	identity interface{},
	f func(accumulator interface{}, element string) interface{},
) interface{} {
	result := identity

	s.ForEach(func(element string) {
		result = f(result, element)
	})

	return result
}

// Skip returns a new StringStream that skips the first n elements
func (s StringStream) Skip(n int) StringStream {
	done := false

	return StringStream{
		iterator: func() (string, bool) {
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

// Sorted returns a new StringStream with the values sorted by the provided comparator..
func (s StringStream) Sorted() StringStream {
	var sortedIter func() (string, bool)
	done := false

	return StringStream{
		iterator: func() (string, bool) {
			if !done {
				// Sort all StringStream elements
				sorted := s.ToSlice()
				sort.Strings(sorted)

				sortedIter = (&stringSliceIterator{array: sorted}).next
				done = true
			}

			// Return next sorted element
			return sortedIter()
		},
	}
}

// ToMap returns a map of all elements by invoking the given function to a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s StringStream) ToMap(f func(string) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	s.ForEach(func(element string) {
		k, v := f(element)
		m[k] = v
	})

	return m
}

// ToSlice returns a slice of all elements
func (s StringStream) ToSlice() []string {
	var array []string

	s.ForEach(func(element string) {
		array = append(array, element)
	})

	return array
}
