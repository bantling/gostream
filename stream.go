package gostream

import (
	"reflect"
	"sort"
	"sync"

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

const (
	// DefaultNumberOfParallelItems is the default number of items when executing transforms in parallel
	DefaultNumberOfParallelItems uint = 50
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

// compose two func(Iter) Iter f1, f2 and returns a composition func(x Iter) Iter of f2(f1(x))
// If f1 is nil, the composition degenerates to f2(x)
// Panics if f2 is nil
func compose(f1, f2 func(*goiter.Iter) *goiter.Iter) func(*goiter.Iter) *goiter.Iter {
	if f2 == nil {
		panic("compose: f2 cannot be nil")
	}

	composition := f2
	if f1 != nil {
		composition = func(it *goiter.Iter) *goiter.Iter {
			return f2(f1(it))
		}
	}

	return composition
}

// doParallel does the grunt work of parallel processing, returning a slice of results.
// If numItems is 0, the default value is DefaultNumberOfParallelItems.
func doParallel(
	source *goiter.Iter,
	transform func(*goiter.Iter) *goiter.Iter,
	finisher func(*goiter.Iter) *goiter.Iter,
	numItems uint,
	flag ParallelFlags,
) []interface{} {
	n := DefaultNumberOfParallelItems
	if numItems > 0 {
		n = numItems
	}

	var flatData []interface{}
	if transform == nil {
		// If the transform is nil, there is no transform, just use source vales as is
		flatData = source.ToSlice()
	} else {
		var splitData [][]interface{}
		if flag == NumberOfGoroutines {
			// numItems = desired number of rows; number of colums to be determined
			splitData = source.SplitIntoColumns(n)
		} else {
			// numItems = desired number of columns; number of rows to be determined
			splitData = source.SplitIntoRows(n)
		}

		// Execute goroutines, one per row of splitData.
		// Each goroutine applies the queued operations to each item in its row.
		wg := &sync.WaitGroup{}

		for i, row := range splitData {
			wg.Add(1)

			go func(i int, row []interface{}) {
				defer wg.Done()

				splitData[i] = transform(goiter.OfElements(row)).ToSlice()
			}(i, row)
		}

		// Wait for all goroutines to complete
		wg.Wait()

		// Combine rows into a single flat slice
		flatData = goiter.FlattenArraySlice(splitData)
	}

	// If the finisher is non-nil, apply it afterwards - it cannot be done in parallel
	if finisher != nil {
		flatData = finisher(goiter.Of(flatData...)).ToSlice()
	}

	// Return transformed rows
	return flatData
}

// Stream is based on a source iterator, and provides a streaming facility where items can be transformed one by one as they are iterated into a new set, and possibly apply further transforms on the new set.
// A Stream is effectively a kind of builder pattern, building up a set of transforms from an input data set to an output data set.
//
// The idea is to compose a set of transforms, then call a terminal method that will invoke the composed transforms and produce a new result.
// All single element transforms are handled by Stream.
// All multi element transforms are handled by Finisher.
//
// The Transform method allow for arbitrary transforms, for cases where the transforms provided are not sufficient.
// When calling the transform methods, the transforms are composed using function composition so that there is only one transform function in the Stream.
// Each transform is a function that acceps an Iter and returns a new Iter.
//
// For example, suppose the following sequence is executed:
//
// Stream.
//   Given(1,3,1,2,9,7,2,4,7,5,8,6,8).
//   Filter(FilterFunc(func(i int) bool { return i < 5 })).
//   Map(MapFunc(func(i int) int { return i * 2 })).
//   AndThen().
//   Distinct().
//   Sort(IntSortFunc).
//   ToSliceOf(0)
//
// The order of operations is exactly as indicated - filter then map each element one by one into a new set, finally remove duplicates from then sort the set.
// The result will be []int{2,4,6,8}.
type Stream struct {
	source    *goiter.Iter
	transform func(*goiter.Iter) *goiter.Iter
	finite    bool
}

// construct handles the details common to all constructor functions
func construct(source *goiter.Iter, finite bool) Stream {
	return Stream{
		source:    source,
		transform: nil,
		finite:    finite,
	}
}

// ==== Constructors

// Of constructs a stream of hard-coded values
func Of(items ...interface{}) Stream {
	return construct(
		goiter.Of(items...),
		true,
	)
}

// OfIterables constructs a stream of values returned by any number of iterables
func OfIterables(iterables ...goiter.Iterable) Stream {
	return construct(
		goiter.OfIterables(iterables...),
		true,
	)
}

// Iterate returns a stream of an infinite iterative calculation, f(seed), f(f(seed)), ...
// Since the series is infinite, some combination of Stream.First() and/or Finisher.Limit() will be required to terminate the series.
func Iterate(seed interface{}, f func(interface{}) interface{}) Stream {
	acculumator := seed

	return construct(
		goiter.NewIter(func() (interface{}, bool) {
			acculumator = f(acculumator)

			return acculumator, true
		}),
		false,
	)
}

// === Transforms

// Transform composes the current transform with a new one
func (s Stream) Transform(t func(*goiter.Iter) *goiter.Iter) Stream {
	return Stream{
		source:    s.source,
		transform: compose(s.transform, t),
		finite:    s.finite,
	}
}

// Filter returns a new stream of all elements that pass the given predicate
func (s Stream) Filter(f func(element interface{}) bool) Stream {
	return s.Transform(
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

// Map maps each element to a new element, possibly of a different type
func (s Stream) Map(f func(element interface{}) interface{}) Stream {
	return s.Transform(
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
	return s.Transform(
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

// Iter returns an iterator of the elements in this Stream.
// Note that a stream can only be iterated once by a single *goiter.Iter instance.
func (s Stream) Iter() *goiter.Iter {
	it := s.source
	if s.transform != nil {
		it = s.transform(it)
	}

	return goiter.NewIter(
		func() (interface{}, bool) {
			if it.Next() {
				return it.Value(), true
			}

			return nil, false
		},
	)
}

// AndThen returns a Finisher, which performs additional post processing on the results of the transforms in this Stream.
func (s Stream) AndThen() Finisher {
	return Finisher{
		source:    s,
		transform: nil,
	}
}

// FindFirst returns the optional first element of applying any tranforms to the Stream.
// See Finisher.FindFirst().
func (s Stream) FindFirst() gooptional.Optional {
	return s.AndThen().FindFirst()
}

// ToSlice returns a slice of the transformed elements.
// See Finisher.ToSlice().
func (s Stream) ToSlice() []interface{} {
	return s.AndThen().ToSlice()
}

// ToSliceOf returns a slice of the transformed elements.
// See Finisher.ToSliceOf().
func (s Stream) ToSliceOf(val interface{}) interface{} {
	return s.AndThen().ToSliceOf(val)
}

// ==== Finisher

// Finisher does two things:
// 1. Apply zero or more transforms that operate across multiple elements after any Stream transforms have been applied to each individual element of the Stream source
// 2. Provide terminal methods that return the final result of applying the Stream and Finisher trasforms to the Stream source
//
// The purpose of separating Finisher from Stream is twofold:
// 1. Make the chaining method calls accurately represent that all multi-element transforms are applied after all single element tranforms.
// 2. Simplify paralell execution of transforms by breaking it into two phases:
//    a. Execute single element transforms on the Stream source in parallel
//    b. Execute multi element transforms on the result of the parallel execution
//
// Guaranteeing the mutli element transforms occur after parallel execution of single element transforms greatly simplifies the parallel algorithm:
// - Only one parallel algorithm is needed
// - No need for multiple passes or buffering
//
// If the Stream was constructed using the Iterate constructor function, then all terminal Finisher methods will panic unless the Limit(int) method is called first.
// This guarantees no infinite loop will occur in the terminal methods.
type Finisher struct {
	source    Stream
	transform func(*goiter.Iter) *goiter.Iter
}

// FindFirst returns the optional first element of applying any tranforms to the stream source.
// May be called any number of times at any time.
// Exhausts one or more items of the source until an item that satisfies the current transforms is found, if any.
// If no such item is found, an empty Optional is returned, else an Optional of the transformed item is returned.
//
// Note that it is possible for the transforms to transform an item into a nil value, resulting in an empty Optional.
// A such, an empty result does not necessarily indicate there are no more results in the Stream.
// However, Stream is based on goiter.Iter, which panics if Next() is called again after a previous Next() call returned false.
// Taken together, the FindFirst() result cannot distinguish between a nil element and the end of the stream.
func (fin Finisher) FindFirst() gooptional.Optional {
	var val interface{}

	it := fin.source.Iter()
	if fin.transform != nil {
		it = fin.transform(it)
	}

	if it.Next() {
		val = it.Value()
	}

	return gooptional.Of(val)
}

// ==== Transforms

// Transform composes the current transform with a new one
func (fin Finisher) Transform(f func(*goiter.Iter) *goiter.Iter) Finisher {
	return Finisher{
		source:    fin.source,
		transform: compose(fin.transform, f),
	}
}

// Filter returns a new Finisher of all elements that pass the given predicate
func (fin Finisher) Filter(f func(element interface{}) bool) Finisher {
	return fin.Transform(
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
func (fin Finisher) FilterNot(f func(element interface{}) bool) Finisher {
	return fin.Filter(
		func(element interface{}) bool {
			return !f(element)
		},
	)
}

// Distinct returns a Finisher of distinct elements only
func (fin Finisher) Distinct() Finisher {
	alreadyRead := map[interface{}]bool{}

	return fin.Filter(
		func(element interface{}) bool {
			if !alreadyRead[element] {
				alreadyRead[element] = true
				return true
			}

			return false
		},
	)
}

// Duplicates returns a stream of duplicate elements only
func (fin Finisher) Duplicates() Finisher {
	alreadyRead := map[interface{}]bool{}

	return fin.Filter(
		func(element interface{}) bool {
			if !alreadyRead[element] {
				alreadyRead[element] = true
				return false
			}

			return true
		},
	)
}

// Skip returns a new stream that skips the first n elements
func (fin Finisher) Skip(n int) Finisher {
	skipped := false

	return fin.Transform(
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

// Limit returns a new stream that only iterates the first n elements, ignoring the rest
func (fin Finisher) Limit(n uint) Finisher {
	var (
		elementsRead uint
	)

	return fin.Transform(
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

// Sorted returns a new stream with the values sorted by the provided comparator.
func (fin Finisher) Sorted(less func(element1, element2 interface{}) bool) Finisher {
	var sortedIter *goiter.Iter
	done := false

	return fin.Transform(
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
func (fin Finisher) ReverseSorted(less func(element1, element2 interface{}) bool) Finisher {
	return fin.Sorted(func(element1, element2 interface{}) bool {
		return !less(element1, element2)
	})
}

// ==== Terminals

// Iter returns an iterator of the elements in this Finisher.
// Note that a Finisher can only be iterated once by a single *goiter.Iter instance.
func (fin Finisher) Iter() *goiter.Iter {
	it := fin.source.Iter()
	if fin.transform != nil {
		it = fin.transform(it)
	}

	return goiter.NewIter(
		func() (interface{}, bool) {
			if it.Next() {
				return it.Value(), true
			}

			return nil, false
		},
	)
}

// AllMatch is true if the predicate matches all elements with short-circuit logic
func (fin Finisher) AllMatch(f func(element interface{}) bool) bool {
	allMatch := true

	for it := fin.Iter(); it.Next(); {
		if allMatch = f(it.Value()); !allMatch {
			break
		}
	}

	return allMatch
}

// AnyMatch is true if the predicate matches any element with short-circuit logic
func (fin Finisher) AnyMatch(f func(element interface{}) bool) bool {
	anyMatch := false

	for it := fin.Iter(); it.Next(); {
		if anyMatch = f(it.Value()); anyMatch {
			break
		}
	}

	return anyMatch
}

// NoneMatch is true if the predicate matches none of the elements with short-circuit logic
func (fin Finisher) NoneMatch(f func(element interface{}) bool) bool {
	noneMatch := true

	for it := fin.Iter(); it.Next(); {
		if noneMatch = !f(it.Value()); !noneMatch {
			break
		}
	}

	return noneMatch
}

// Average returns an optional average value.
// The slice elements must be convertible to a float64.
func (fin Finisher) Average() gooptional.Optional {
	var (
		sum   float64
		count int
	)

	for it := fin.Iter(); it.Next(); {
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
func (fin Finisher) Sum() gooptional.Optional {
	var (
		sum    float64
		hasSum bool
	)

	for it := fin.Iter(); it.Next(); {
		sum += it.Float64Value()
		hasSum = true
	}

	if !hasSum {
		return gooptional.Of()
	}

	return gooptional.Of(sum)
}

// Count returns the count of all elements
func (fin Finisher) Count() int {
	count := 0

	for it := fin.Iter(); it.Next(); {
		count++
	}

	return count
}

// Last returns the optional last element
func (fin Finisher) Last() gooptional.Optional {
	var last interface{}

	for it := fin.Iter(); it.Next(); {
		last = it.Value()
	}

	return gooptional.Of(last)
}

// Max returns an optional maximum value according to the provided comparator
func (fin Finisher) Max(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var max interface{}

	if it := fin.Iter(); it.Next() {
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
func (fin Finisher) Min(less func(element1, element2 interface{}) bool) gooptional.Optional {
	var min interface{}

	if it := fin.Iter(); it.Next() {
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

// ForEach invokes a consumer with each element of the stream
func (fin Finisher) ForEach(f func(element interface{})) {
	for it := fin.Iter(); it.Next(); {
		f(it.Value())
	}
}

// Reduce uses a function to reduce the stream to a single value by iteratively executing a function
// with the current accumulated value and the next stream element.
// The identity provided is the initial accumulated value, which means the result type is the
// same type as the initial value, which can be any type.
// If there are no elements in the strea, the result is the identity.
// Otherwise, the result is f(f(identity, element1), element2)...
func (fin Finisher) Reduce(
	identity interface{},
	f func(accumulator interface{}, element2 interface{}) interface{},
) interface{} {
	result := identity

	for it := fin.Iter(); it.Next(); {
		result = f(result, it.Value())
	}

	return result
}

// GroupBy groups elements by executing the given function on each value to get a key,
// and appending the element to the end of a slice associated with the key in the resulting map.
func (fin Finisher) GroupBy(f func(element interface{}) (key interface{})) map[interface{}][]interface{} {
	m := map[interface{}][]interface{}{}

	fin.Reduce(
		m,
		func(accumulator interface{}, element interface{}) interface{} {
			k := f(element)
			m[k] = append(m[k], element)
			return m
		},
	)

	return m
}

// ToMap returns a map of all elements by invoking the given function to get a key/value pair for the map.
// It is up to the function to generate unique keys to prevent values from being overwritten.
func (fin Finisher) ToMap(f func(interface{}) (key interface{}, value interface{})) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	for it := fin.Iter(); it.Next(); {
		k, v := f(it.Value())
		m[k] = v
	}

	return m
}

// ToMapOf returns a map of all elements, where the map key and value types are the same as the types of aKey and aValue.
// EG, if aKey is an int and aVaue is a string, then a map[int]string is returned.
// Panics if keys are not convertible to the key type or values are not convertible to the value type.
func (fin Finisher) ToMapOf(
	f func(interface{}) (key interface{}, value interface{}),
	aKey, aValue interface{},
) interface{} {
	var (
		ktyp = reflect.TypeOf(aKey)
		vtyp = reflect.TypeOf(aValue)
		m    = reflect.MakeMap(reflect.MapOf(ktyp, vtyp))
	)

	for it := fin.Iter(); it.Next(); {
		k, v := f(it.Value())
		m.SetMapIndex(
			reflect.ValueOf(k).Convert(ktyp),
			reflect.ValueOf(v).Convert(vtyp),
		)
	}

	return m.Interface()
}

// ToSlice returns a slice of all elements
func (fin Finisher) ToSlice() []interface{} {
	array := []interface{}{}

	for it := fin.Iter(); it.Next(); {
		array = append(array, it.Value())
	}

	return array
}

// ToSliceOf returns a slice of all elements, where the slice elements are the same type as the type of elementVal.
// EG, if elementVal is an int, an []int is returned.
// Panics if elements are not convertible to the type of elementVal.
func (fin Finisher) ToSliceOf(elementVal interface{}) interface{} {
	var (
		elementTyp = reflect.TypeOf(elementVal)
		array      = reflect.MakeSlice(reflect.SliceOf(elementTyp), 0, 0)
	)

	for it := fin.Iter(); it.Next(); {
		array = reflect.Append(array, reflect.ValueOf(it.Value()).Convert(elementTyp))
	}

	return array.Interface()
}

// ToStream returns a stream of all elements
func (fin Finisher) ToStream() Stream {
	return Of(fin.ToSlice()...)
}

// ==== Parallel processing

// ParallelToStream processes the result of the current Finisher in parallel using a number of goroutines.
// The number of items provided is interpreted according to the optional ParallelFlags value:
// 1. NumberOfGoroutines - numItems indicates the number of go routines (default)
// 2. NumberOfItemsPerGoroutine - numItems indicates the number of items each go routine processes
// Either way, the results are ordered, and a new Stream is returned that iterates them.
// If numItems is 0, it defaults to DefaultNumberOfParallelItems.
func (fin Finisher) ParallelToStream(numItems uint, flag ...ParallelFlags) Stream {
	theFlag := NumberOfGoroutines
	if len(flag) > 0 {
		theFlag = flag[0]
	}

	data := doParallel(
		fin.source.source,
		fin.source.transform,
		fin.transform,
		numItems,
		theFlag,
	)

	return Of(data...)
}

// ParallelToSlice is the same as Parallel, except that it returns the data as a slice
func (fin Finisher) ParallelToSlice(numItems uint, flag ...ParallelFlags) []interface{} {
	theFlag := NumberOfGoroutines
	if len(flag) > 0 {
		theFlag = flag[0]
	}

	data := doParallel(
		fin.source.Iter(),
		fin.transform,
		fin.transform,
		numItems,
		theFlag,
	)

	return data
}

// ParallelToSliceOf is the same as ParallelSlice, except that it returns the data as a slice whose type matches the element value given.
func (fin Finisher) ParallelToSliceOf(elementValue interface{}, numItems uint, flag ...ParallelFlags) interface{} {
	theFlag := NumberOfGoroutines
	if len(flag) > 0 {
		theFlag = flag[0]
	}

	data := doParallel(
		fin.source.source,
		fin.source.transform,
		fin.transform,
		numItems,
		theFlag,
	)

	return goiter.FlattenArraySliceAsType(data, elementValue)
}
