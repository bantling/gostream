package gostream

import (
	"strconv"
	"testing"

	"github.com/bantling/goiter"
	"github.com/stretchr/testify/assert"
)

// ==== Constructors

func TestOf(t *testing.T) {
	s := Of(3, 2, 1)
	assert.Equal(t, []interface{}{3, 2, 1}, s.ToSlice())
}

func TestOfIter(t *testing.T) {
	s := OfIter(goiter.OfElements([]int{6, 5, 4}))
	assert.Equal(t, []interface{}{6, 5, 4}, s.ToSlice())
}

func TestStreamIterate(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return element.(int) * 2
	}
	s := Iterate(1, fn)
	first := s.First()
	assert.Equal(t, 2, first.MustGet())
	first = s.First()
	assert.Equal(t, 4, first.MustGet())
	first = s.First()
	assert.Equal(t, 8, first.MustGet())

	fn2 := IterateFunc(func(element int) int {
		return element * 2
	})
	s = Iterate(1, fn2)
	first = s.First()
	assert.Equal(t, 2, first.MustGet())
	first = s.First()
	assert.Equal(t, 4, first.MustGet())
	first = s.First()
	assert.Equal(t, 8, first.MustGet())
}

// ==== Other

func TestStreamFirst(t *testing.T) {
	s := Of()
	first := s.First()
	assert.True(t, first.IsEmpty())

	s = Of(1)
	first = s.First()
	assert.Equal(t, 1, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())

	s = Of(1, 2)
	first = s.First()
	assert.Equal(t, 1, first.MustGet())
	first = s.First()
	assert.Equal(t, 2, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())
}

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
	assert.Equal(t, []interface{}{}, s.Distinct().ToSlice())

	s = Of(1, 1)
	assert.Equal(t, []interface{}{1}, s.Distinct().ToSlice())

	s = Of(1, 2, 2, 1)
	assert.Equal(t, []interface{}{1, 2}, s.Distinct().ToSlice())
}

func TestStreamDuplicate(t *testing.T) {
	s := Of()
	assert.Equal(t, []interface{}{}, s.Duplicate().ToSlice())

	s = Of(1, 1, 2)
	assert.Equal(t, []interface{}{1}, s.Duplicate().ToSlice())

	s = Of(1, 2, 2, 1, 3)
	assert.Equal(t, []interface{}{2, 1}, s.Duplicate().ToSlice())
}

func TestStreamFilter(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.Equal(t, []interface{}{}, s.Filter(fn).ToSlice())

	s = Of(1, 2, 3)
	assert.Equal(t, []interface{}{1, 2}, s.Filter(fn).ToSlice())

	fn2 := FilterFunc(func(element int) bool { return element < 3 })
	s = Of(1, 2, 3)
	assert.Equal(t, []int{1, 2}, s.Filter(fn2).ToSliceOf(0))
}

func TestStreamFilterNot(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.Equal(t, []interface{}{}, s.FilterNot(fn).ToSlice())

	s = Of(1, 2, 3)
	assert.Equal(t, []interface{}{3}, s.FilterNot(fn).ToSlice())
}

func TestStreamLimit(t *testing.T) {
	s := Of(1, 2, 3)
	assert.Equal(t, []interface{}{1, 2}, s.Limit(2).ToSlice())
}

func TestStreamMap(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return strconv.Itoa(element.(int) * 2)
	}
	s := Of().Map(fn)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).Map(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = Of(1, 2).Map(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())

	fn2 := MapFunc(func(element int) string { return strconv.Itoa(element * 2) })
	s = Of(1, 2).Map(fn2)
	assert.Equal(t, []string{"2", "4"}, s.ToSliceOf(""))
}

func TestStreamPeek(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := Of().Peek(fn)
	assert.Equal(t, elements, []interface{}(nil), s.ToSlice())

	elements = nil
	s = Of(1).Peek(fn)
	assert.Equal(t, elements, []interface{}{1}, s.ToSlice())

	elements = nil
	s = Of(1, 2).Peek(fn)
	assert.Equal(t, elements, []interface{}{1, 2}, s.ToSlice())

	var elements2 []int
	fn2 := PeekFunc(func(element int) { elements2 = append(elements2, element) })
	s = Of(1, 2).Peek(fn2)
	assert.Equal(t, elements2, []int{1, 2}, s.ToSliceOf(0))
}

func TestStreamSkip(t *testing.T) {
	s := Of().Skip(0)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).Skip(0)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(1).Skip(1)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1, 2).Skip(1)
	assert.Equal(t, []interface{}{2}, s.ToSlice())

	s = Of(1, 2, 3).Skip(2)
	assert.Equal(t, []interface{}{3}, s.ToSlice())

	s = Of(1, 2, 3, 4).Skip(2)
	assert.Equal(t, []interface{}{3, 4}, s.ToSlice())
}

func TestStreamSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of().Sorted(fn)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).Sorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(2, 1).Sorted(fn)
	assert.Equal(t, []interface{}{1, 2}, s.ToSlice())

	s = Of(2, 3, 1).Sorted(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, s.ToSlice())

	fn2 := func(i, j int) bool { return i < j }
	s = Of(2, 1).Sorted(SortFunc(fn2))
	assert.Equal(t, []int{1, 2}, s.ToSliceOf(0))
}

func TestStreamReverseSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of().ReverseSorted(fn)
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = Of(2, 1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{2, 1}, s.ToSlice())

	s = Of(2, 3, 1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{3, 2, 1}, s.ToSlice())

	fn2 := func(i, j int) bool { return i < j }
	s = Of(1, 2).ReverseSorted(SortFunc(fn2))
	assert.Equal(t, []int{2, 1}, s.ToSliceOf(0))
}

// ==== Terminals

func TestStreamAllMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.True(t, s.AllMatch(fn))

	s = Of(1, 2)
	assert.True(t, s.AllMatch(fn))

	s = Of(1, 2, 3)
	assert.False(t, s.AllMatch(fn))

	s = Of(1, 2, 3, 4)
	assert.False(t, s.AllMatch(fn))
}

func TestStreamAnyMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.False(t, s.AnyMatch(fn))

	s = Of(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = Of(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestStreamAverage(t *testing.T) {
	s := Of(1, 2.25)
	avg := (1 + 2.25) / 2
	assert.Equal(t, avg, s.Average().Iter().NextFloatValue())
}

func TestStreamSum(t *testing.T) {
	s := Of(1, 2.25)
	sum := 1 + 2.25
	assert.Equal(t, sum, s.Sum().Iter().NextFloatValue())
}

func TestStreamNoneMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := Of()
	assert.True(t, s.NoneMatch(fn))

	s = Of(3, 4)
	assert.True(t, s.NoneMatch(fn))

	s = Of(1, 2, 3)
	assert.False(t, s.NoneMatch(fn))

	s = Of(1, 2, 3, 4)
	assert.False(t, s.NoneMatch(fn))
}

func TestStreamCount(t *testing.T) {
	s := Of()
	assert.Equal(t, 0, s.Count())

	s = Of(2)
	assert.Equal(t, 1, s.Count())

	s = Of(2, 3)
	assert.Equal(t, 2, s.Count())
}

func TestStreamForEach(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := Of()
	s.ForEach(fn)
	assert.Equal(t, []interface{}(nil), elements)

	elements = nil
	s = Of(1)
	s.ForEach(fn)
	assert.Equal(t, []interface{}{1}, elements)

	elements = nil
	s = Of(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, elements)
}

func TestStreamGroupBy(t *testing.T) {
	fn := func(element interface{}) (key interface{}) {
		return element.(int) % 3
	}
	s := Of()
	assert.Equal(t, map[interface{}][]interface{}{}, s.GroupBy(fn))

	s = Of(0)
	assert.Equal(t, map[interface{}][]interface{}{0: {0}}, s.GroupBy(fn))

	s = Of(0, 1, 4)
	assert.Equal(t, map[interface{}][]interface{}{0: {0}, 1: {1, 4}}, s.GroupBy(fn))
}

func TestStreamLast(t *testing.T) {
	s := Of()
	last := s.Last()
	assert.True(t, last.IsEmpty())

	s = Of(1)
	last = s.Last()
	assert.Equal(t, 1, last.MustGet())

	s = Of(1, 2)
	last = s.Last()
	assert.Equal(t, 2, last.MustGet())
}

func TestStreamMax(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of()
	max := s.Max(fn)
	assert.True(t, max.IsEmpty())

	s = Of(1)
	max = s.Max(fn)
	assert.Equal(t, 1, max.MustGet())

	s = Of(1, 2)
	max = s.Max(fn)
	assert.Equal(t, 2, max.MustGet())

	s = Of(1, 3, 2)
	max = s.Max(fn)
	assert.Equal(t, 3, max.MustGet())
}

func TestStreamMin(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := Of()
	min := s.Min(fn)
	assert.True(t, min.IsEmpty())

	s = Of(1)
	min = s.Min(fn)
	assert.Equal(t, 1, min.MustGet())

	s = Of(1, 0)
	min = s.Min(fn)
	assert.Equal(t, 0, min.MustGet())

	s = Of(1, -1, 2)
	min = s.Min(fn)
	assert.Equal(t, -1, min.MustGet())
}

func TestStreamReduce(t *testing.T) {
	fn := func(accumulator, element2 interface{}) interface{} {
		return accumulator.(int) + element2.(int)
	}
	s := Of()
	sum := s.Reduce(0, fn)
	assert.Equal(t, 0, sum)

	s = Of(1, 2, 3)
	sum = s.Reduce(1, fn)
	assert.Equal(t, 7, sum)
}

func TestStreamToMap(t *testing.T) {
	fn := func(element interface{}) (k interface{}, v interface{}) {
		return element, strconv.Itoa(element.(int))
	}
	s := Of()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = Of(1)
	assert.Equal(t, map[interface{}]interface{}{1: "1"}, s.ToMap(fn))

	s = Of(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1: "1", 2: "2", 3: "3"}, s.ToMap(fn))
}

func TestStreamToSlice(t *testing.T) {
	s := Of()
	assert.Equal(t, []interface{}{}, s.ToSlice())

	s = Of(1, 2)
	assert.Equal(t, []interface{}{1, 2}, s.ToSlice())
}

func TestStreamToSliceOf(t *testing.T) {
	s := Of()
	assert.Equal(t, []int{}, s.ToSliceOf(0))

	s = Of(1, 2)
	assert.Equal(t, []int{1, 2}, s.ToSliceOf(0))
}
