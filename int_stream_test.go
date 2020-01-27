package gostream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSliceIterator(t *testing.T) {
	ai := intSliceIterator{array: []int{1, 2, 3}}
	next, hasNext := ai.next()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, 2, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, 3, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.False(t, hasNext)
}

func TestIntStreamNewIntStream(t *testing.T) {
	ai := intSliceIterator{array: []int{1, 2, 3}}
	s := NewIntStream(ai.next)
	assert.Equal(t, []int{1, 2, 3}, s.ToSlice())

	s = NewIntStreamOf(3, 2, 1)
	assert.Equal(t, []int{3, 2, 1}, s.ToSlice())
}

func TestIntStreamAllMatch(t *testing.T) {
	fn := func(element int) bool { return element < 3 }
	s := NewIntStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = NewIntStreamOf(1, 2)
	assert.True(t, s.AllMatch(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.False(t, s.AllMatch(fn))

	s = NewIntStreamOf(1, 2, 3, 4)
	assert.False(t, s.AllMatch(fn))
}

func TestIntStreamAnyMatch(t *testing.T) {
	fn := func(element int) bool { return element < 3 }
	s := NewIntStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = NewIntStreamOf(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestIntStreamAverage(t *testing.T) {
	s := NewIntStreamOf()
	avg := s.Average()
	assert.True(t, avg.IsEmpty())

	s = NewIntStreamOf(3, 4)
	avg = s.Average()
	assert.Equal(t, 3.5, avg.MustGet())
	avg = s.Average()
	assert.True(t, avg.IsEmpty())

	s = NewIntStreamOf(1, 2, 3)
	avg = s.Average()
	assert.Equal(t, 2.0, avg.MustGet())
	avg = s.Average()
	assert.True(t, avg.IsEmpty())
}

func TestIntStreamConcat(t *testing.T) {
	s1 := NewIntStreamOf(1, 2, 3)
	s2 := NewIntStreamOf(4, 5, 6)
	s3 := s1.Concat(s2)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s3.ToSlice())
}

func TestIntStreamCount(t *testing.T) {
	s := NewIntStreamOf()
	assert.Equal(t, 0, s.Count())

	s = NewIntStreamOf(2)
	assert.Equal(t, 1, s.Count())

	s = NewIntStreamOf(2, 3)
	assert.Equal(t, 2, s.Count())
}

func TestIntStreamDistinct(t *testing.T) {
	s := NewIntStreamOf()
	assert.Equal(t, []int(nil), s.Distinct().ToSlice())

	s = NewIntStreamOf(1, 1)
	assert.Equal(t, []int{1}, s.Distinct().ToSlice())

	s = NewIntStreamOf(1, 2, 2, 1)
	assert.Equal(t, []int{1, 2}, s.Distinct().ToSlice())
}

func TestIntStreamDuplicate(t *testing.T) {
	s := NewIntStreamOf()
	assert.Equal(t, []int(nil), s.Duplicate().ToSlice())

	s = NewIntStreamOf(1, 1, 2)
	assert.Equal(t, []int{1}, s.Duplicate().ToSlice())

	s = NewIntStreamOf(1, 2, 2, 1, 3)
	assert.Equal(t, []int{2, 1}, s.Duplicate().ToSlice())
}

func TestIntStreamFilter(t *testing.T) {
	fn := func(element int) bool { return element < 3 }
	s := NewIntStreamOf()
	assert.Equal(t, []int(nil), s.Filter(fn).ToSlice())

	s = NewIntStreamOf(1, 2, 3)
	assert.Equal(t, []int{1, 2}, s.Filter(fn).ToSlice())
}

func TestIntStreamFirst(t *testing.T) {
	s := NewIntStreamOf()
	first := s.First()
	assert.True(t, first.IsEmpty())

	s = NewIntStreamOf(1)
	first = s.First()
	assert.Equal(t, 1, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())

	s = NewIntStreamOf(1, 2)
	first = s.First()
	assert.Equal(t, 1, first.MustGet())
	first = s.First()
	assert.Equal(t, 2, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())
}

func TestIntStreamForEach(t *testing.T) {
	var elements []int
	fn := func(element int) {
		elements = append(elements, element)
	}
	s := NewIntStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []int(nil), elements)

	elements = nil
	s = NewIntStreamOf(1)
	s.ForEach(fn)
	assert.Equal(t, []int{1}, elements)

	elements = nil
	s = NewIntStreamOf(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []int{1, 2, 3}, elements)
}

func TestIntStreamGroupBy(t *testing.T) {
	fn := func(element int) (key interface{}) {
		return element % 3
	}
	s := NewIntStreamOf()
	assert.Equal(t, map[interface{}][]int{}, s.GroupBy(fn))

	s = NewIntStreamOf(0)
	assert.Equal(t, map[interface{}][]int{0: []int{0}}, s.GroupBy(fn))

	s = NewIntStreamOf(0, 1, 4)
	assert.Equal(t, map[interface{}][]int{0: []int{0}, 1: []int{1, 4}}, s.GroupBy(fn))
}

func TestIntStreamIterate(t *testing.T) {
	fn := func(element int) int {
		return element * 2
	}
	s := NewIntStreamOf().Iterate(1, fn)
	next := s.First()
	assert.Equal(t, 2, next.MustGet())
	next = s.First()
	assert.Equal(t, 4, next.MustGet())
	next = s.First()
	assert.Equal(t, 8, next.MustGet())
}

func TestIntStreamLast(t *testing.T) {
	s := NewIntStreamOf()
	last := s.Last()
	assert.True(t, last.IsEmpty())

	s = NewIntStreamOf(1)
	last = s.Last()
	assert.Equal(t, 1, last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())

	s = NewIntStreamOf(1, 2)
	last = s.Last()
	assert.Equal(t, 2, last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())
}

func TestIntStreamMap(t *testing.T) {
	fn := func(element int) int {
		return element * 2
	}
	s := NewIntStreamOf().Map(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewIntStreamOf(1).Map(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = NewIntStreamOf(1, 2).Map(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestIntStreamMapToFloat(t *testing.T) {
	fn := func(element int) float64 {
		return float64(element * 2)
	}
	s := NewIntStreamOf().MapToFloat(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewIntStreamOf(1).MapToFloat(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = NewIntStreamOf(1, 2).MapToFloat(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestIntStreamMapTo(t *testing.T) {
	fn := func(element int) interface{} {
		return element * 2
	}
	s := NewIntStreamOf().MapTo(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewIntStreamOf(1).MapTo(fn)
	assert.Equal(t, []interface{}{2}, s.ToSlice())

	s = NewIntStreamOf(1, 2).MapTo(fn)
	assert.Equal(t, []interface{}{2, 4}, s.ToSlice())
}

func TestIntStreamMapToString(t *testing.T) {
	fn := func(element int) string {
		return strconv.Itoa(element * 2)
	}
	s := NewIntStreamOf().MapToString(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewIntStreamOf(1).MapToString(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = NewIntStreamOf(1, 2).MapToString(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestIntStreamMax(t *testing.T) {
	s := NewIntStreamOf()
	max := s.Max()
	assert.True(t, max.IsEmpty())

	s = NewIntStreamOf(1)
	max = s.Max()
	assert.Equal(t, 1, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = NewIntStreamOf(1, 2)
	max = s.Max()
	assert.Equal(t, 2, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = NewIntStreamOf(1, 3, 2)
	max = s.Max()
	assert.Equal(t, 3, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())
}

func TestIntStreamMin(t *testing.T) {
	s := NewIntStreamOf()
	min := s.Min()
	assert.True(t, min.IsEmpty())

	s = NewIntStreamOf(1)
	min = s.Min()
	assert.Equal(t, 1, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = NewIntStreamOf(1, 0)
	min = s.Min()
	assert.Equal(t, 0, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = NewIntStreamOf(1, -1, 2)
	min = s.Min()
	assert.Equal(t, -1, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())
}

func TestIntStreamNoneMatch(t *testing.T) {
	fn := func(element int) bool { return element < 3 }
	s := NewIntStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = NewIntStreamOf(3, 4)
	assert.True(t, s.NoneMatch(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.False(t, s.NoneMatch(fn))

	s = NewIntStreamOf(1, 2, 3, 4)
	assert.False(t, s.NoneMatch(fn))
}

func TestIntStreamPeek(t *testing.T) {
	var elements []int
	fn := func(element int) {
		elements = append(elements, element)
	}
	s := NewIntStreamOf().Peek(fn)
	assert.Equal(t, elements, []int(nil), s.ToSlice())

	elements = nil
	s = NewIntStreamOf(1).Peek(fn)
	assert.Equal(t, elements, []int{1}, s.ToSlice())

	elements = nil
	s = NewIntStreamOf(1, 2).Peek(fn)
	assert.Equal(t, elements, []int{1, 2}, s.ToSlice())
}

func TestIntStreamReduce(t *testing.T) {
	fn := func(accumulator interface{}, element int) interface{} {
		return accumulator.(int) + element
	}
	s := NewIntStreamOf()
	sum := s.Reduce(0, fn)
	assert.Equal(t, 0, sum)

	s = NewIntStreamOf(1, 2, 3)
	sum = s.Reduce(1, fn)
	assert.Equal(t, 7, sum)
}

func TestIntStreamReverseSorted(t *testing.T) {
	s := NewIntStreamOf().ReverseSorted()
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewIntStreamOf(1).ReverseSorted()
	assert.Equal(t, []int{1}, s.ToSlice())

	s = NewIntStreamOf(1, 2).ReverseSorted()
	assert.Equal(t, []int{2, 1}, s.ToSlice())

	s = NewIntStreamOf(2, 3, 1).ReverseSorted()
	assert.Equal(t, []int{3, 2, 1}, s.ToSlice())
}

func TestIntStreamSkip(t *testing.T) {
	s := NewIntStreamOf().Skip(0)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewIntStreamOf(1).Skip(0)
	assert.Equal(t, []int{1}, s.ToSlice())

	s = NewIntStreamOf(1).Skip(1)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewIntStreamOf(1, 2).Skip(1)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = NewIntStreamOf(1, 2, 3).Skip(2)
	assert.Equal(t, []int{3}, s.ToSlice())

	s = NewIntStreamOf(1, 2, 3, 4).Skip(2)
	assert.Equal(t, []int{3, 4}, s.ToSlice())
}

func TestIntStreamSorted(t *testing.T) {
	s := NewIntStreamOf().Sorted()
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewIntStreamOf(1).Sorted()
	assert.Equal(t, []int{1}, s.ToSlice())

	s = NewIntStreamOf(2, 1).Sorted()
	assert.Equal(t, []int{1, 2}, s.ToSlice())

	s = NewIntStreamOf(2, 3, 1).Sorted()
	assert.Equal(t, []int{1, 2, 3}, s.ToSlice())
}

func TestIntStreamSum(t *testing.T) {
	s := NewIntStreamOf()
	sum := s.Sum()
	assert.True(t, sum.IsEmpty())

	s = NewIntStreamOf(3, 4)
	sum = s.Sum()
	assert.Equal(t, 7, sum.MustGet())
	sum = s.Sum()
	assert.True(t, sum.IsEmpty())

	s = NewIntStreamOf(1, 2, 3)
	sum = s.Sum()
	assert.Equal(t, 6, sum.MustGet())
	sum = s.Sum()
	assert.True(t, sum.IsEmpty())
}

func TestIntStreamToMap(t *testing.T) {
	fn := func(element int) (k interface{}, v interface{}) {
		return element, strconv.Itoa(element)
	}
	s := NewIntStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = NewIntStreamOf(1)
	assert.Equal(t, map[interface{}]interface{}{1: "1"}, s.ToMap(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1: "1", 2: "2", 3: "3"}, s.ToMap(fn))
}
