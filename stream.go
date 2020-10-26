package gostream

import (
	"math/cmplx"
	"reflect"
	"sort"
	//	"sync"

	"github.com/bantling/goiter"
	"github.com/bantling/gooptional"
)

// ParallelFlags is a pair of flags indicating whether to interpret the number as the number of goroutines or the number of items each goroutine processes
type ParallelFlags uint

const (
	// NumberOfGoroutines is the default, and indicates the number of goroutines
	NumberOfGoroutines ParallelFlags = iota
	// NumberOfItemsPerGoroutine indicates the number of items each goroutine processes
	NumberOfItemsPerGoroutine
)

// IterateFunc adapts any func that accepts and returns the exact same type into func(interface{}) interface{} suitable for the Iterate method.
// Panics if f is not a func that accepts and returns one type that is exactly the same.
func IterateFunc(f interface{}) func(interface{}) interface{} {
	var (
		val = reflect.ValueOf(f)
		typ = val.Type()
	)

	if typ.Kind() != reflect.Func {
		panic("f must be a function")
	}

	if (typ.NumIn() != 1) || (typ.NumOut() != 1) {
		panic("f must accept and return a single value of the exact same type")
	}

	argType, retType := typ.In(0), typ.Out(0)
	if argType != retType {
		panic("f must accept and return a single value of the exact same type")
	}

	return func(arg interface{}) interface{} {
		return val.Call([]reflect.Value{reflect.ValueOf(arg)})[0].Interface()
	}
}

// FilterFunc adapts any func that accepts a single arg and returns bool into a func(interface{}) bool suitable for the Filter methods.
// Panics if f is not a func that accepts a single arg and returns bool.
var FilterFunc = gooptional.FilterFunc

// MapFunc adapts any func that accepts a single arg and returns a single value into a func(interface{}) interface{} suitable for the Map method.
// Panics if f is not a func that accepts a single arg and returns a single value.
var MapFunc = gooptional.MapFunc

// PeekFunc adapts any func that accepts a single arg and returns nothing into a func(interface{}) suitable for the Peek method.
// Panics if f is not a func that accepts a single arg and returns nothing.
var PeekFunc = gooptional.ConsumerFunc

// SortFunc adapts any func that accepts a pair of args of exactly the same type and returns true if first arg < second arg, suitable for the Sort method.
// Panics if f is not a func that accepts a pair of args of exactly the same type and returns bool.
func SortFunc(f interface{}) func(i, j interface{}) bool {
	var (
		val = reflect.ValueOf(f)
		typ = val.Type()
	)

	if typ.Kind() != reflect.Func {
		panic("f must be a function")
	}

	if (typ.NumIn() != 2) || (typ.NumOut() != 1) {
		panic("f must accept a pair of args of exactly the same type and returns bool")
	}

	if (typ.In(0) != typ.In(1)) || (typ.Out(0).Kind() != reflect.Bool) {
		panic("f must accept a pair of args of exactly the same type and returns bool")
	}

	return func(i, j interface{}) bool {
		return val.Call([]reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)})[0].Bool()
	}
}

// IntSortFunc sorts int values
func IntSortFunc(i, j interface{}) bool {
	return i.(int) < j.(int)
}

// Int8SortFunc sorts int8 values
func Int8SortFunc(i, j interface{}) bool {
	return i.(int8) < j.(int8)
}

// Int16SortFunc sorts int16 values
func Int16SortFunc(i, j interface{}) bool {
	return i.(int16) < j.(int16)
}

// Int32SortFunc sorts int32 values
func Int32SortFunc(i, j interface{}) bool {
	return i.(int32) < j.(int32)
}

// Int64SortFunc sorts int64 values
func Int64SortFunc(i, j interface{}) bool {
	return i.(int64) < j.(int64)
}

// UintSortFunc sorts uint values
func UintSortFunc(i, j interface{}) bool {
	return i.(uint) < j.(uint)
}

// Uint8SortFunc sorts uint8 values
func Uint8SortFunc(i, j interface{}) bool {
	return i.(uint8) < j.(uint8)
}

// Uint16SortFunc sorts uint16 values
func Uint16SortFunc(i, j interface{}) bool {
	return i.(uint16) < j.(uint16)
}

// Uint32SortFunc sorts uint32 values
func Uint32SortFunc(i, j interface{}) bool {
	return i.(uint32) < j.(uint32)
}

// Uint64SortFunc sorts uint64 values
func Uint64SortFunc(i, j interface{}) bool {
	return i.(uint64) < j.(uint64)
}

// Float32SortFunc sorts float32 values
func Float32SortFunc(i, j interface{}) bool {
	return i.(float32) < j.(float32)
}

// Float64SortFunc sorts float64 values
func Float64SortFunc(i, j interface{}) bool {
	return i.(float64) < j.(float64)
}

// Complex64SortFunc sorts complex64 values
func Complex64SortFunc(i, j interface{}) bool {
	return cmplx.Abs(complex128(i.(complex64))) < cmplx.Abs(complex128(j.(complex64)))
}

// Complex128SortFunc sorts complex128 values
func Complex128SortFunc(i, j interface{}) bool {
	return cmplx.Abs(i.(complex128)) < cmplx.Abs(j.(complex128))
}

// StringSortFunc sorts string values
func StringSortFunc(i, j interface{}) bool {
	return i.(string) < j.(string)
}

// Stream is based on an iterator, and provides a streaming facility where items can be filtered and transformed as they are iterated.
// The idea is to call some non-terminal methods to queue up a set of operations, then call a terminal method that will invoke the queued up operations and produce a new result.
// All non-terminal methods return a new Stream.
// Terminal methods either return a value or invoke a function for every element.
// There are also parallel processing methods, which do the following:
// 1. Split up the original set of items provided when the Stream was constructed into rows of sub slices.
// 2. Use a separate go routine for each row to process the queued up operations on each element, producing a new row.
// 3. Return the new rows as is, as a single row, or as a new Stream.
// The paralell processing methods are all terminal, as they all iterate all the items provided when the Stream was constructed.
// The ParallelAsStream method that returns a new Stream is also a non-terminal operation because further method chaining can continue to queue up further operations on the results of the paralell processing.
type Stream struct {
	source *goiter.Iter
	queue  func(*goiter.Iter) *goiter.Iter
}

// ==== Helpers

// construct handles the details common to all constructor functions
func construct(source *goiter.Iter) Stream {
	return Stream{
		source: source,
		queue:  nil,
	}
}

// addQueue handles the details of adding another function to the queue for this stream
func (s Stream) addQueue(f func(*goiter.Iter) *goiter.Iter) Stream {
	nq := f
	if s.queue != nil {
		nq = func(it *goiter.Iter) *goiter.Iter {
			return f(s.queue(it))
		}
	}

	return Stream{
		source: s.source,
		queue:  nq,
	}
}

// doParallel does the grunt work of parallel processing, returning a slice of results
//func (s *Stream) doParallel(numItems uint, f ParallelFlags) []interface{} {
//	var splitData [][]interface{}
//	if f == NumberOfGoroutines {
//		// numItems = desired number of rows, number of colums to be determined
//		splitData = s.queue.Iter().SplitIntoColumns(numItems)
//	} else {
//		// numItems = desired number of columns. number of rows to be determined
//		splitData = s.queue.Iter().SplitIntoRows(numItems)
//	}
//
//	// Execute goroutines, one per row of splitData.
//	// Each goroutines applies the queued operations for each item in its row.
//	wg := &sync.WaitGroup{}
//	for i, row := range splitData {
//		wg.Add(1)
//
//		go func(wg *sync.WaitGroup, row []interface{}, splitData [][]interface{}, i int) {
//			defer wg.Done()
//
//			s := construct(goiter.OfElements(row))
//			splitData[i] = s.ToSlice()
//		}(wg, row, splitData, i)
//	}
//	wg.Wait()
//
//	// After all goroutines have completed, combine results into a single slice
//	return goiter.FlattenArraySlice(splitData)
//}

// ==== Constructors

// Of constructs a stream of hard-coded values
func Of(items ...interface{}) Stream {
	return construct(goiter.Of(items...))
}

// OfIterables constructs a stream of values returned by any number of iterables
func OfIterables(iterables ...goiter.Iterable) Stream {
	return construct(goiter.OfIterables(iterables...))
}

// Iterate returns a stream of an iterative calculation, f(seed), f(f(seed)), ...
func Iterate(seed interface{}, f func(interface{}) interface{}) Stream {
	acculumator := seed

	return construct(
		goiter.NewIter(func() (interface{}, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		}),
	)
}

// ==== Other

// Iter is the goiter.Iterable interface, returns an iterator of the results of this stream
func (s Stream) Iter() *goiter.Iter {
	if s.queue == nil {
		return s.source
	}

	return s.queue(s.source)
}

// First returns the optional first element
func (s Stream) First() gooptional.Optional {
	var val interface{}

	if it := s.Iter(); it.Next() {
		val = it.Value()
	}

	return gooptional.Of(val)
}

// ==== Filters

// Filter returns a new stream of all elements that pass the given predicate
func (s Stream) Filter(f func(element interface{}) bool) Stream {
	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					for it.Next() {
						if val := it.Value(); f(val) {
							return val, true
						}
					}

					return nil, false
				},
			)
		},
	)
}

// FilterNot returns a new stream of all elements that do not pass the given predicate
func (s Stream) FilterNot(f func(element interface{}) bool) Stream {
	return s.Filter(
		func(element interface{}) bool {
			return !f(element)
		},
	)
}

// Distinct returns a stream of distinct elements only
func (s Stream) Distinct() Stream {
	alreadyRead := map[interface{}]bool{}

	return s.Filter(
		func(element interface{}) bool {
			if !alreadyRead[element] {
				alreadyRead[element] = true
				return true
			}

			return false
		},
	)
}

// Duplicate returns a stream of duplicate elements only
func (s Stream) Duplicate() Stream {
	alreadyRead := map[interface{}]bool{}

	return s.Filter(
		func(element interface{}) bool {
			if !alreadyRead[element] {
				alreadyRead[element] = true
				return false
			}

			return true
		},
	)
}

// Limit returns a new stream that only iterates the first n elements, ignoring the rest
func (s Stream) Limit(n uint) Stream {
	var (
		elementsRead uint
	)

	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					if (elementsRead == n) || (!it.Next()) {
						return nil, false
					}

					elementsRead++
					return it.Value(), true
				},
			)
		},
	)
}

// Map maps each element to a new element, possibly of a different type
func (s Stream) Map(f func(element interface{}) interface{}) Stream {
	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					if it.Next() {
						return f(it.Value()), true
					}

					return nil, false
				},
			)
		},
	)
}

// Peek returns a stream that calls a function that examines each value and performs an additional operation
func (s Stream) Peek(f func(interface{})) Stream {
	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					if it.Next() {
						val := it.Value()
						f(val)
						return val, true
					}

					return nil, false
				},
			)
		},
	)
}

// Skip returns a new stream that skips the first n elements
func (s Stream) Skip(n int) Stream {
	skipped := false

	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					// Skip n elements only once
					if !skipped {
						skipped = true

						for i := 1; i <= n; i++ {
							if !it.Next() {
								// We don't have n elements to skip
								return nil, false
							}
						}
					}

					if it.Next() {
						// Return next element
						return it.Value(), true
					}

					return nil, false
				},
			)
		},
	)
}

// Sorted returns a new stream with the values sorted by the provided comparator..
func (s Stream) Sorted(less func(element1, element2 interface{}) bool) Stream {
	var sortedIter *goiter.Iter
	done := false

	return s.addQueue(
		func(it *goiter.Iter) *goiter.Iter {
			return goiter.NewIter(
				func() (interface{}, bool) {
					if !done {
						// Sort all stream elements
						sorted := it.ToSlice()
						sort.Slice(sorted, func(i, j int) bool {
							return less(sorted[i], sorted[j])
						})

						sortedIter = goiter.OfElements(sorted)
						done = true
					}

					// Return next sorted element
					if sortedIter.Next() {
						return sortedIter.Value(), true
					}

					return nil, false
				},
			)
		},
	)
}

// ReverseSorted returns a stream with elements sorted in decreasing order.
// The provided function must compare elements in increasing order, same as for Sorted.
func (s Stream) ReverseSorted(less func(element1, element2 interface{}) bool) Stream {
	return s.Sorted(func(element1, element2 interface{}) bool {
		return !less(element1, element2)
	})
}

// ==== Terminals

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (s Stream) AllMatch(f func(element interface{}) bool) bool {
	allMatch := true

	for it := s.Iter(); it.Next(); {
		if allMatch = f(it.Value()); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s Stream) AnyMatch(f func(element interface{}) bool) bool {
	anyMatch := false

	for it := s.Iter(); it.Next(); {
		if anyMatch = f(it.Value()); anyMatch {
			break
		}
	}

	return anyMatch
}

// Average returns an optional average value.
// The slice elements must be convertible to a float64.
func (s Stream) Average() gooptional.Optional {
	var (
		sum   float64
		count int
	)

	for it := s.Iter(); it.Next(); {
		sum += it.Float64Value()
		count++
	}

	if count == 0 {
		return gooptional.Of()
	}

	avg := sum / float64(count)
	return gooptional.Of(avg)
}

// Sum returns an optional sum value.
// The slice elements must be convertible to a float64.
func (s Stream) Sum() gooptional.Optional {
	var (
		sum    float64
		hasSum bool
	)

	for it := s.Iter(); it.Next(); {
		sum += it.Float64Value()
		hasSum = true
	}

	if !hasSum {
		return gooptional.Of()
	}

	return gooptional.Of(sum)
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (s Stream) NoneMatch(f func(element interface{}) bool) bool {
	noneMatch := true

	for it := s.Iter(); it.Next(); {
		if noneMatch = !f(it.Value()); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Count returns the count of all elements
func (s Stream) Count() int {
	count := 0

	for it := s.Iter(); it.Next(); {
		count++
	}

	return count
}

// ForEach invokes a consumer with each element of the stream
func (s Stream) ForEach(f func(element interface{})) {
	for it := s.Iter(); it.Next(); {
		f(it.Value())
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

// Last returns the optional last element
func (s Stream) Last() gooptional.Optional {
	var last interface{}

	for it := s.Iter(); it.Next(); {
		last = it.Value()
	}

	return gooptional.Of(last)
}

// Max returns an optional maximum value according to the provided comparator
func (s Stream) Max(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var max interface{}

	if it := s.Iter(); it.Next() {
		max = it.Value()

		for it.Next() {
			element := it.Value()

			if less(max, element) {
				max = element
			}
		}
	}

	return gooptional.Of(max)
}

// Min returns an optional minimum value according to the provided comparator
func (s Stream) Min(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var (
		min interface{}
		it  = s.Iter()
	)

	if it.Next() {
		min = it.Value()

		for it.Next() {
			element := it.Value()

			if less(element, min) {
				min = element
			}
		}
	}

	return gooptional.Of(min)
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

	for it := s.Iter(); it.Next(); {
		result = f(result, it.Value())
	}

	return result
}

// ToMap returns a map of all elements by invoking the given function to get a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s Stream) ToMap(f func(interface{}) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	for it := s.Iter(); it.Next(); {
		k, v := f(it.Value())
		m[k] = v
	}

	return m
}

// ToSlice returns a slice of all elements
func (s Stream) ToSlice() []interface{} {
	array := []interface{}{}

	for it := s.Iter(); it.Next(); {
		array = append(array, it.Value())
	}

	return array
}

// ToSliceOf returns a slice of all elements, where the slice type is the same as the given element.
// EG, if a value of type int is passed, an []int is returned.
func (s Stream) ToSliceOf(elementVal interface{}) interface{} {
	array := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(elementVal)), 0, 0)

	for it := s.Iter(); it.Next(); {
		array = reflect.Append(array, reflect.ValueOf(it.Value()))
	}

	return array.Interface()
}

// ==== Paralell processing
