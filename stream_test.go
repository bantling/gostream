// SPDX-License-Identifier: Apache-2.0

package gostream

import (
	//	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/bantling/gofuncs"
	"github.com/bantling/goiter"
	"github.com/stretchr/testify/assert"
)

// ==== Constructors

func TestOf(t *testing.T) {
	s := Of(3, 2, 1)
	assert.Equal(t, []interface{}{3, 2, 1}, s.AndThen().ToSlice())
}

func TestOfIterables(t *testing.T) {
	s := OfIterables(goiter.OfElements([]int{6, 5, 4}))
	assert.Equal(t, []interface{}{6, 5, 4}, s.AndThen().ToSlice())
}

func TestStreamIterate(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return element.(int) * 2
	}
	fin := Iterate(1, fn).AndThen()
	first := fin.FindFirst()
	assert.Equal(t, 2, first.MustGet())
	first = fin.FindFirst()
	assert.Equal(t, 4, first.MustGet())
	first = fin.FindFirst()
	assert.Equal(t, 8, first.MustGet())

	fn2 := IterateFunc(func(element int) int {
		return element * 2
	})
	fin = Iterate(1, fn2).AndThen()
	first = fin.FindFirst()
	assert.Equal(t, 2, first.MustGet())
	first = fin.FindFirst()
	assert.Equal(t, 4, first.MustGet())
	first = fin.FindFirst()
	assert.Equal(t, 8, first.MustGet())

	// Panic on infinite Finisher
	func() {
		defer func() {
			assert.Equal(t, ErrInfiniteFinisher, recover())
		}()

		fin.ToSlice()
		assert.Fail(t, "Must panic")
	}()

	// Apply limit to make it finite
	assert.Equal(t, []int{16, 32, 64, 128}, fin.Limit(4).ToSliceOf(0))
}

// ==== Other

func TestStreamIsIterable(t *testing.T) {
	var (
		s                        = Of(1)
		iterable goiter.Iterable = s
		it       *goiter.Iter    = iterable.Iter()
	)

	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())
	assert.False(t, it.Next())
}

// ==== Transforms

func TestStreamDistinct(t *testing.T) {
	s := Of()
	assert.Equal(t, []interface{}{}, s.AndThen().Distinct().ToSlice())

	s = Of(1, 1)
	assert.Equal(t, []interface{}{1}, s.AndThen().Distinct().ToSlice())

	s = Of(1, 2, 2, 1)
	assert.Equal(t, []interface{}{1, 2}, s.AndThen().Distinct().ToSlice())
}

func TestStreamDuplicates(t *testing.T) {
	s := Of()
	assert.Equal(t, []interface{}{}, s.AndThen().Duplicates().ToSlice())

	s = Of(1, 1, 2)
	assert.Equal(t, []interface{}{1}, s.AndThen().Duplicates().ToSlice())

	s = Of(1, 2, 2, 1, 3)
	assert.Equal(t, []interface{}{2, 1}, s.AndThen().Duplicates().ToSlice())
}

func TestStreamFilter(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.Equal(t, []interface{}{}, s.Filter(fn).AndThen().ToSlice())

	s = Of(1, 2, 3)
	assert.Equal(t, []interface{}{1, 2}, s.Filter(fn).AndThen().ToSlice())

	fn2 := gofuncs.Filter(func(element int) bool { return element < 3 })
	s = Of(1, 2, 3)
	assert.Equal(t, []int{1, 2}, s.Filter(fn2).AndThen().ToSliceOf(0))
}

func TestStreamFilterNot(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.Equal(t, []interface{}{}, s.FilterNot(fn).AndThen().ToSlice())

	s = Of(1, 2, 3)
	assert.Equal(t, []interface{}{3}, s.FilterNot(fn).AndThen().ToSlice())
}

func TestStreamLimit(t *testing.T) {
	s := Of(1, 2, 3)
	assert.Equal(t, []interface{}{1, 2}, s.AndThen().Limit(2).ToSlice())
}

func TestStreamMap(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return strconv.Itoa(element.(int) * 2)
	}
	s := Of().Map(fn)
	assert.Equal(t, []interface{}{}, s.AndThen().ToSlice())

	s = Of(1).Map(fn)
	assert.Equal(t, []interface{}{"2"}, s.AndThen().ToSlice())

	s = Of(1, 2).Map(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.AndThen().ToSlice())

	fn2 := gofuncs.Map(func(element int) string { return strconv.Itoa(element * 2) })
	s = Of(1, 2).Map(fn2)
	assert.Equal(t, []string{"2", "4"}, s.AndThen().ToSliceOf(""))
}

func TestStreamPeek(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := Of().Peek(fn)
	assert.Equal(t, elements, []interface{}(nil), s.AndThen().ToSlice())

	elements = nil
	s = Of(1).Peek(fn)
	assert.Equal(t, elements, []interface{}{1}, s.AndThen().ToSlice())

	elements = nil
	s = Of(1, 2).Peek(fn)
	assert.Equal(t, elements, []interface{}{1, 2}, s.AndThen().ToSlice())

	var elements2 []int
	fn2 := gofuncs.Consumer(func(element int) { elements2 = append(elements2, element) })
	s = Of(1, 2).Peek(fn2)
	assert.Equal(t, elements2, []int{1, 2}, s.AndThen().ToSliceOf(0))
}

func TestStreamSkip(t *testing.T) {
	s := Of().AndThen().Skip(0)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).AndThen().Skip(0)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(1).AndThen().Skip(1)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1, 2).AndThen().Skip(1)
	assert.Equal(t, []interface{}{2}, s.ToSlice())

	s = Of(1, 2, 3).AndThen().Skip(2)
	assert.Equal(t, []interface{}{3}, s.ToSlice())

	s = Of(1, 2, 3, 4).AndThen().Skip(2)
	assert.Equal(t, []interface{}{3, 4}, s.ToSlice())
}

func TestStreamSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of().AndThen().Sorted(fn)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).AndThen().Sorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(2, 1).AndThen().Sorted(fn)
	assert.Equal(t, []interface{}{1, 2}, s.ToSlice())

	s = Of(2, 3, 1).AndThen().Sorted(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, s.ToSlice())

	fn2 := func(i, j int) bool { return i < j }
	s = Of(2, 1).AndThen().Sorted(gofuncs.SortFunc(fn2))
	assert.Equal(t, []int{1, 2}, s.ToSliceOf(0))

	s = Of(2, 1).AndThen().Sorted(gofuncs.SortFunc(fn2))
	assert.Equal(t, []uint{1, 2}, s.ToSliceOf(uint(0)))

	s = Of(2, 1).AndThen().Sorted(gofuncs.IntSortFunc)
	assert.Equal(t, []int{1, 2}, s.ToSliceOf(0))

	s = Of(int8(2), int8(1)).AndThen().Sorted(gofuncs.IntSortFunc)
	assert.Equal(t, []int8{1, 2}, s.ToSliceOf(int8(0)))

	s = Of(uint(2), uint(1)).AndThen().Sorted(gofuncs.UintSortFunc)
	assert.Equal(t, []uint{1, 2}, s.ToSliceOf(uint(0)))

	s = Of(uint8(2), uint8(1)).AndThen().Sorted(gofuncs.UintSortFunc)
	assert.Equal(t, []uint8{1, 2}, s.ToSliceOf(uint8(0)))

	s = Of(float32(2), float32(1)).AndThen().Sorted(gofuncs.FloatSortFunc)
	assert.Equal(t, []float32{1, 2}, s.ToSliceOf(float32(0)))

	s = Of(float64(2), float64(1)).AndThen().Sorted(gofuncs.FloatSortFunc)
	assert.Equal(t, []float64{1, 2}, s.ToSliceOf(0.0))

	s = Of(complex64(1+2i), complex64(2+3i)).AndThen().Sorted(gofuncs.ComplexSortFunc)
	assert.Equal(t, []complex64{(1 + 2i), (2 + 3i)}, s.ToSliceOf(complex64(0)))

	s = Of(complex128(2+3i), complex128(1+2i)).AndThen().Sorted(gofuncs.ComplexSortFunc)
	assert.Equal(t, []complex128{(1 + 2i), (2 + 3i)}, s.ToSliceOf(complex128(0)))

	s = Of("b", "a").AndThen().Sorted(gofuncs.StringSortFunc)
	assert.Equal(t, []string{"a", "b"}, s.ToSliceOf(""))

	s = Of('b', 'a').AndThen().Sorted(gofuncs.StringSortFunc)
	assert.Equal(t, []string{"a", "b"}, s.ToSliceOf(""))

	s = Of(big.NewInt(2), big.NewInt(1)).AndThen().Sorted(gofuncs.BigIntSortFunc)
	assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2)}, s.ToSliceOf((*big.Int)(nil)))

	s = Of(big.NewRat(2, 3), big.NewRat(1, 2)).AndThen().Sorted(gofuncs.BigRatSortFunc)
	assert.Equal(t, []*big.Rat{big.NewRat(1, 2), big.NewRat(2, 3)}, s.ToSliceOf((*big.Rat)(nil)))

	s = Of(big.NewFloat(2.0), big.NewFloat(1.0)).AndThen().Sorted(gofuncs.BigFloatSortFunc)
	assert.Equal(t, []*big.Float{big.NewFloat(1.0), big.NewFloat(2.0)}, s.ToSliceOf((*big.Float)(nil)))
}

func TestStreamReverseSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of().AndThen().ReverseSorted(fn)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).AndThen().ReverseSorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(2, 1).AndThen().ReverseSorted(fn)
	assert.Equal(t, []interface{}{2, 1}, s.ToSlice())

	s = Of(2, 3, 1).AndThen().ReverseSorted(fn)
	assert.Equal(t, []interface{}{3, 2, 1}, s.ToSlice())

	fn2 := func(i, j int) bool { return i < j }
	s = Of(1, 2).AndThen().ReverseSorted(gofuncs.SortFunc(fn2))
	assert.Equal(t, []int{2, 1}, s.ToSliceOf(0))
}

// ==== Terminals

func TestStreamAllMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.True(t, s.AndThen().AllMatch(fn))

	s = Of(1, 2)
	assert.True(t, s.AndThen().AllMatch(fn))

	s = Of(1, 2, 3)
	assert.False(t, s.AndThen().AllMatch(fn))

	s = Of(1, 2, 3, 4)
	assert.False(t, s.AndThen().AllMatch(fn))
}

func TestStreamAnyMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.False(t, s.AndThen().AnyMatch(fn))

	s = Of(3, 4)
	assert.False(t, s.AndThen().AnyMatch(fn))

	s = Of(1, 2, 3)
	assert.True(t, s.AndThen().AnyMatch(fn))
}

func TestStreamAverage(t *testing.T) {
	s := Of(1, 2.25)
	avg := (1 + 2.25) / 2
	assert.Equal(t, avg, s.AndThen().Average().Iter().NextFloat64Value())
}

func TestStreamSum(t *testing.T) {
	s := Of(1, 2.25)
	sum := 1 + 2.25
	assert.Equal(t, sum, s.AndThen().Sum().Iter().NextFloat64Value())
}

func TestStreamNoneMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.True(t, s.AndThen().NoneMatch(fn))

	s = Of(3, 4)
	assert.True(t, s.AndThen().NoneMatch(fn))

	s = Of(1, 2, 3)
	assert.False(t, s.AndThen().NoneMatch(fn))

	s = Of(1, 2, 3, 4)
	assert.False(t, s.AndThen().NoneMatch(fn))
}

func TestStreamCount(t *testing.T) {
	s := Of()
	assert.Equal(t, 0, s.AndThen().Count())

	s = Of(2)
	assert.Equal(t, 1, s.AndThen().Count())

	s = Of(2, 3)
	assert.Equal(t, 2, s.AndThen().Count())
}

func TestStreamForEach(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := Of()
	s.AndThen().ForEach(fn)
	assert.Equal(t, []interface{}(nil), elements)

	elements = nil
	s = Of(1)
	s.AndThen().ForEach(fn)
	assert.Equal(t, []interface{}{1}, elements)

	elements = nil
	s = Of(1, 2, 3)
	s.AndThen().ForEach(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, elements)
}

func TestStreamGroupBy(t *testing.T) {
	fn := func(element interface{}) (key interface{}) {
		return element.(int) % 3
	}
	s := Of()
	assert.Equal(t, map[interface{}][]interface{}{}, s.AndThen().GroupBy(fn))

	s = Of(0)
	assert.Equal(t, map[interface{}][]interface{}{0: {0}}, s.AndThen().GroupBy(fn))

	s = Of(0, 1, 4)
	assert.Equal(t, map[interface{}][]interface{}{0: {0}, 1: {1, 4}}, s.AndThen().GroupBy(fn))
}

func TestStreamLast(t *testing.T) {
	s := Of()
	last := s.AndThen().Last()
	assert.True(t, last.IsEmpty())

	s = Of(1)
	last = s.AndThen().Last()
	assert.Equal(t, 1, last.MustGet())

	s = Of(1, 2)
	last = s.AndThen().Last()
	assert.Equal(t, 2, last.MustGet())
}

func TestStreamMax(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of()
	max := s.AndThen().Max(fn)
	assert.True(t, max.IsEmpty())

	s = Of(1)
	max = s.AndThen().Max(fn)
	assert.Equal(t, 1, max.MustGet())

	s = Of(1, 2)
	max = s.AndThen().Max(fn)
	assert.Equal(t, 2, max.MustGet())

	s = Of(1, 3, 2)
	max = s.AndThen().Max(fn)
	assert.Equal(t, 3, max.MustGet())
}

func TestStreamMin(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of()
	min := s.AndThen().Min(fn)
	assert.True(t, min.IsEmpty())

	s = Of(1)
	min = s.AndThen().Min(fn)
	assert.Equal(t, 1, min.MustGet())

	s = Of(1, 0)
	min = s.AndThen().Min(fn)
	assert.Equal(t, 0, min.MustGet())

	s = Of(1, -1, 2)
	min = s.AndThen().Min(fn)
	assert.Equal(t, -1, min.MustGet())
}

func TestStreamReduce(t *testing.T) {
	fn := func(accumulator, element2 interface{}) interface{} {
		return accumulator.(int) + element2.(int)
	}
	s := Of()
	sum := s.AndThen().Reduce(0, fn)
	assert.Equal(t, 0, sum)

	s = Of(1, 2, 3)
	sum = s.AndThen().Reduce(1, fn)
	assert.Equal(t, 7, sum)
}

func TestStreamToMap(t *testing.T) {
	fn := func(element interface{}) (k interface{}, v interface{}) {
		return element, strconv.Itoa(element.(int))
	}
	s := Of()
	assert.Equal(t, map[interface{}]interface{}{}, s.AndThen().ToMap(fn))

	s = Of(1)
	assert.Equal(t, map[interface{}]interface{}{1: "1"}, s.AndThen().ToMap(fn))

	s = Of(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1: "1", 2: "2", 3: "3"}, s.AndThen().ToMap(fn))
}

func TestStreamToSlice(t *testing.T) {
	s := Of()
	assert.Equal(t, []interface{}{}, s.AndThen().ToSlice())

	s = Of(1, 2)
	assert.Equal(t, []interface{}{1, 2}, s.AndThen().ToSlice())

	s = Of(1, 2).Filter(gofuncs.Filter(func(i int) bool { return i > 5 }))
	assert.Equal(t, []interface{}{}, s.AndThen().ToSlice())
}

func TestStreamToSliceOf(t *testing.T) {
	s := Of()
	assert.Equal(t, []int{}, s.AndThen().ToSliceOf(0))

	s = Of(1, 2)
	assert.Equal(t, []int{1, 2}, s.AndThen().ToSliceOf(0))
}

// ==== Sequence

func TestSequence(t *testing.T) {
	s := Of(1, 2, 1, 3, 4, 3, 5, 6, 7, 7, 8, 9, 10).
		Map(gofuncs.Map(func(i int) int { return i * 2 })).
		//  2, 4,  2, 6, 8, 6, 10, 12, 14, 14, 16, 18, 20
		Map(gofuncs.Map(func(i int) int { return i - 3 })).
		// -1, 1, -1, 3, 5, 3,  7,  9, 11, 11, 13, 15, 17
		Filter(gofuncs.Filter(func(i int) bool { return i <= 7 })).
		// -1, 1, -1, 3, 5, 3,  7
		AndThen().
		Skip(2).
		//  -1, 3, 5, 3, 7
		Distinct().
		// -1, 3, 5, 7
		ReverseSorted(gofuncs.IntSortFunc).
		//  7, 5,  3, -1
		Limit(3)
		// 7, 5, 3
	assert.Equal(t, 7, s.FindFirst().MustGet())
	// 5, 3
	assert.Equal(t, []int{5, 3}, s.ToSliceOf(0))
}

func TestParallel(t *testing.T) {
	var (
		doubler         = gofuncs.Map(func(i int) int { return i * 2 })
		input           = []int{1, 2, 1, 3, 4, 3, 5, 6, 7, 7, 8, 9, 10}
		empty           = []int{}
		doubled         = []int{2, 4, 2, 6, 8, 6, 10, 12, 14, 14, 16, 18, 20}
		distinct        = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		doubledDistinct = []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	)

	// a series of tests for all 8 combinations
	// elements?, transform?, finisher?

	// 000
	s := Of().AndThen()
	assert.Equal(t, empty, s.ParallelToSliceOf(0, 0))

	// 001
	s = Of().AndThen().Distinct()
	assert.Equal(t, empty, s.ParallelToSliceOf(0, 0))

	// 010
	s = Of().Map(doubler).AndThen()
	assert.Equal(t, empty, s.ParallelToSliceOf(0, 0))

	// 011
	s = Of().Map(doubler).AndThen().Distinct()
	assert.Equal(t, empty, s.ParallelToSliceOf(0, 0))

	// 100
	s = OfIterables(goiter.OfElements(input)).AndThen()
	assert.Equal(t, input, s.ParallelToSliceOf(0, 0))

	// 101
	s = OfIterables(goiter.OfElements(input)).AndThen().Distinct()
	assert.Equal(t, distinct, s.ParallelToSliceOf(0, 0))

	// 110
	s = OfIterables(goiter.OfElements(input)).Map(doubler).AndThen()
	assert.Equal(t, doubled, s.ParallelToSliceOf(0, 0))

	// 111
	s = OfIterables(goiter.OfElements(input)).Map(doubler).AndThen().Distinct()
	assert.Equal(t, doubledDistinct, s.ParallelToSliceOf(0, 0))
}
