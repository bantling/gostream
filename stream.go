package gostream

import (
	"math/cmplx"
	"reflect"
	"sort"

	"github.com/bantling/goiter"
	"github.com/bantling/gooptional"
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
func FilterFunc(f interface{}) func(interface{}) bool {
	var (
		val = reflect.ValueOf(f)
		typ = val.Type()
	)

	if typ.Kind() != reflect.Func {
		panic("f must be a function")
	}

	if (typ.NumIn() != 1) || (typ.NumOut() != 1) {
		panic("f must accept a single arg and return bool")
	}

	if typ.Out(0).Kind() != reflect.Bool {
		panic("f must accept a single arg and return bool")
	}

	return func(arg interface{}) bool {
		return val.Call([]reflect.Value{reflect.ValueOf(arg)})[0].Bool()
	}
}

// MapFunc adapts any func that accepts a single arg and returns a single value into a func(interface{}) interface{} suitable for the Map method.
// Panics if f is not a func that accepts a single arg and returns a single value.
func MapFunc(f interface{}) func(interface{}) interface{} {
	var (
		val = reflect.ValueOf(f)
		typ = val.Type()
	)

	if typ.Kind() != reflect.Func {
		panic("f must be a function")
	}

	if (typ.NumIn() != 1) || (typ.NumOut() != 1) {
		panic("f must accept a single arg and return a single value")
	}

	return func(arg interface{}) interface{} {
		return val.Call([]reflect.Value{reflect.ValueOf(arg)})[0].Interface()
	}
}

// PeekFunc adapts any func that accepts a single arg and returns nothing into a func(interface{}) suitable for the Peek method.
// Panics if f is not a func that accepts a single arg and returns nothing.
func PeekFunc(f interface{}) func(interface{}) {
	var (
		val = reflect.ValueOf(f)
		typ = val.Type()
	)

	if typ.Kind() != reflect.Func {
		panic("f must be a function")
	}

	if (typ.NumIn() != 1) || (typ.NumOut() != 0) {
		panic("f must accept a single arg and return nothing")
	}

	return func(arg interface{}) {
		val.Call([]reflect.Value{reflect.ValueOf(arg)})
	}
}

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

// Stream is the base object type for streams, based on an iterator.
// Some methods are chaining methods, they return a new stream.
// Some functions are terminal, they return a non-stream result.
// Some terminal functions return optional values by returning (<type>, bool), like an iterating function.
// Some functions accept a comparator that takes two elements and returns true if element1 < element 2.
type Stream struct {
	iter *goiter.Iter
}

// ==== Constructors

// Of constructs a stream of hard-coded values
func Of(items ...interface{}) Stream {
	return Stream{iter: goiter.Of(items...)}
}

// OfIter constructs a stream of values returned by an existing iter
func OfIter(iter *goiter.Iter) Stream {
	return Stream{iter: iter}
}

// Iterate returns a stream of an iterative calculation, f(seed), f(f(seed)), ...
func Iterate(seed interface{}, f func(interface{}) interface{}) Stream {
	acculumator := seed

	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		}),
	}
}

// ==== Other

// First returns the optional first element
func (s Stream) First() gooptional.Optional {
	if s.iter.Next() {
		return gooptional.Of(s.iter.Value())
	}

	return gooptional.Of()
}

// Iter is the goiter.Iterable interface, returns an iterator of the elements in this stream
func (s Stream) Iter() *goiter.Iter {
	return s.iter
}

// ==== Transforms

// Distinct returns a stream of distinct elements only
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

// Duplicate returns a stream of duplicate elements only
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
		iter: goiter.NewIter(func() (interface{}, bool) {
			for s.iter.Next() {
				next := s.iter.Value()
				if f(next) {
					return next, true
				}
			}

			return nil, false
		}),
	}
}

// FilterNot returns a new stream of all elements that do not pass the given predicate
func (s Stream) FilterNot(f func(element interface{}) bool) Stream {
	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			for s.iter.Next() {
				next := s.iter.Value()
				if !f(next) {
					return next, true
				}
			}

			return nil, false
		}),
	}
}

// Limit returns a new stream that only iterates the first n elements, ignoring the rest
func (s Stream) Limit(n uint) Stream {
	var (
		elementsRead uint
		done         bool
	)

	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			if done {
				return nil, false
			}

			if !s.iter.Next() {
				done = true
				return nil, false
			}

			elementsRead++
			done = elementsRead == n
			return s.iter.Value(), true
		}),
	}
}

// Map maps each element to a new element, possibly of a different type
func (s Stream) Map(f func(element interface{}) interface{}) Stream {
	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			if s.iter.Next() {
				return f(s.iter.Value()), true
			}

			return nil, false
		}),
	}
}

// Peek returns a stream that calls a function that examines each value and performs an additional operation
func (s Stream) Peek(f func(interface{})) Stream {
	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			if s.iter.Next() {
				val := s.iter.Value()
				f(val)
				return val, true
			}

			return nil, false
		}),
	}
}

// Skip returns a new stream that skips the first n elements
func (s Stream) Skip(n int) Stream {
	var (
		done     = false
		haveMore = true
	)

	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			// Skip n elements only once
			if !done {
				done = true

				for i := 1; i <= n; i++ {
					if !s.iter.Next() {
						haveMore = false
						return nil, false
					}
				}
			}

			if haveMore {
				if haveMore = s.iter.Next(); haveMore {
					// Return next element
					return s.iter.Value(), true
				}
			}

			return nil, false
		}),
	}
}

// Sorted returns a new stream with the values sorted by the provided comparator..
func (s Stream) Sorted(less func(element1, element2 interface{}) bool) Stream {
	var sortedIter *goiter.Iter
	done := false

	return Stream{
		iter: goiter.NewIter(func() (interface{}, bool) {
			if !done {
				// Sort all stream elements
				sorted := s.ToSlice()
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
		}),
	}
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

	for s.iter.Next() {
		if allMatch = f(s.iter.Value()); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (s Stream) AnyMatch(f func(element interface{}) bool) bool {
	anyMatch := false

	for s.iter.Next() {
		if anyMatch = f(s.iter.Value()); anyMatch {
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

	for s.iter.Next() {
		sum += s.iter.FloatValue()
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

	for s.iter.Next() {
		sum += s.iter.FloatValue()
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

	for s.iter.Next() {
		if noneMatch = !f(s.iter.Value()); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Count returns the count of all elements
func (s Stream) Count() int {
	count := 0

	for s.iter.Next() {
		count++
	}

	return count
}

// ForEach invokes a consumer with each element of the stream
func (s Stream) ForEach(f func(element interface{})) {
	for s.iter.Next() {
		f(s.iter.Value())
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
	var (
		last    interface{}
		hasLast bool
	)

	for s.iter.Next() {
		last = s.iter.Value()
		hasLast = true
	}

	if hasLast {
		return gooptional.Of(last)
	}

	return gooptional.Of()
}

// Max returns an optional maximum value according to the provided comparator
func (s Stream) Max(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var max interface{}
	if s.iter.Next() {
		max = s.iter.Value()
		for s.iter.Next() {
			element := s.iter.Value()
			if less(max, element) {
				max = element
			}
		}
	}

	return gooptional.Of(max)
}

// Min returns an optional minimum value according to the provided comparator
func (s Stream) Min(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var min interface{}
	if s.iter.Next() {
		min = s.iter.Value()
		for s.iter.Next() {
			element := s.iter.Value()
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

	for s.iter.Next() {
		result = f(result, s.iter.Value())
	}

	return result
}

// ToMap returns a map of all elements by invoking the given function to get a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (s Stream) ToMap(f func(interface{}) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	for s.iter.Next() {
		k, v := f(s.iter.Value())
		m[k] = v
	}

	return m
}

// ToSlice returns a slice of all elements
func (s Stream) ToSlice() []interface{} {
	array := []interface{}{}

	for s.iter.Next() {
		array = append(array, s.iter.Value())
	}

	return array
}

// ToSliceOf returns a slice of all elements, where the slice type is the same as the given element.
// EG, if a value of type int is passed, an []int is returned.
func (s Stream) ToSliceOf(elementVal interface{}) interface{} {
	array := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(elementVal)), 0, 0)

	for s.iter.Next() {
		array = reflect.Append(array, reflect.ValueOf(s.iter.Value()))
	}

	return array.Interface()
}
