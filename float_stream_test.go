package gostream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatSliceIterator(t *testing.T) {
	ai := floatSliceIterator{array: []float64{1, 2, 3}}
	next, hasNext := ai.next()
	assert.Equal(t, 1.0, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, 2.0, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, 3.0, next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.False(t, hasNext)
}

func TestFloatStreamFloatStreamFromIter(t *testing.T) {
	ai := floatSliceIterator{array: []float64{1, 2, 3}}
	s := FloatStreamFromIter(ai.next)
	assert.Equal(t, []float64{1, 2, 3}, s.ToSlice())

	s = FloatStreamOf(3, 2, 1)
	assert.Equal(t, []float64{3, 2, 1}, s.ToSlice())
}

func TestFloatStreamAllMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := FloatStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = FloatStreamOf(1, 2)
	assert.True(t, s.AllMatch(fn))

	s = FloatStreamOf(1, 2, 3)
	assert.False(t, s.AllMatch(fn))

	s = FloatStreamOf(1, 2, 3, 4)
	assert.False(t, s.AllMatch(fn))
}

func TestFloatStreamAnyMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := FloatStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = FloatStreamOf(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = FloatStreamOf(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestFloatStreamAverage(t *testing.T) {
	s := FloatStreamOf()
	avg := s.Average()
	assert.True(t, avg.IsEmpty())

	s = FloatStreamOf(3, 4)
	avg = s.Average()
	assert.Equal(t, 3.5, avg.MustGet())
	avg = s.Average()
	assert.True(t, avg.IsEmpty())

	s = FloatStreamOf(1, 2, 3)
	avg = s.Average()
	assert.Equal(t, 2.0, avg.MustGet())
	avg = s.Average()
	assert.True(t, avg.IsEmpty())
}

func TestFloatStreamConcat(t *testing.T) {
	s1 := FloatStreamOf(1, 2, 3)
	s2 := FloatStreamOf(4, 5, 6)
	s3 := s1.Concat(s2)
	assert.Equal(t, []float64{1, 2, 3, 4, 5, 6}, s3.ToSlice())
}

func TestFloatStreamCount(t *testing.T) {
	s := FloatStreamOf()
	assert.Equal(t, 0, s.Count())

	s = FloatStreamOf(2)
	assert.Equal(t, 1, s.Count())

	s = FloatStreamOf(2, 3)
	assert.Equal(t, 2, s.Count())
}

func TestFloatStreamDistinct(t *testing.T) {
	s := FloatStreamOf()
	assert.Equal(t, []float64(nil), s.Distinct().ToSlice())

	s = FloatStreamOf(1, 1)
	assert.Equal(t, []float64{1}, s.Distinct().ToSlice())

	s = FloatStreamOf(1, 2, 2, 1)
	assert.Equal(t, []float64{1, 2}, s.Distinct().ToSlice())
}

func TestFloatStreamDuplicate(t *testing.T) {
	s := FloatStreamOf()
	assert.Equal(t, []float64(nil), s.Duplicate().ToSlice())

	s = FloatStreamOf(1, 1, 2)
	assert.Equal(t, []float64{1}, s.Duplicate().ToSlice())

	s = FloatStreamOf(1, 2, 2, 1, 3)
	assert.Equal(t, []float64{2, 1}, s.Duplicate().ToSlice())
}

func TestFloatStreamFilter(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := FloatStreamOf()
	assert.Equal(t, []float64(nil), s.Filter(fn).ToSlice())

	s = FloatStreamOf(1, 2, 3)
	assert.Equal(t, []float64{1, 2}, s.Filter(fn).ToSlice())
}

func TestFloatStreamFirst(t *testing.T) {
	s := FloatStreamOf()

	s = FloatStreamOf(1)
	first := s.First()
	assert.Equal(t, 1.0, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())

	s = FloatStreamOf(1, 2)
	first = s.First()
	assert.Equal(t, 1.0, first.MustGet())
	first = s.First()
	assert.Equal(t, 2.0, first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())
}

func TestFloatStreamForEach(t *testing.T) {
	var elements []float64
	fn := func(element float64) {
		elements = append(elements, element)
	}
	s := FloatStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []float64(nil), elements)

	elements = nil
	s = FloatStreamOf(1)
	s.ForEach(fn)
	assert.Equal(t, []float64{1}, elements)

	elements = nil
	s = FloatStreamOf(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []float64{1, 2, 3}, elements)
}

func TestFloatStreamGroupBy(t *testing.T) {
	fn := func(element float64) (key interface{}) {
		return int(element) % 3
	}
	s := FloatStreamOf()
	assert.Equal(t, map[interface{}][]float64{}, s.GroupBy(fn))

	s = FloatStreamOf(0)
	assert.Equal(t, map[interface{}][]float64{0: {0}}, s.GroupBy(fn))

	s = FloatStreamOf(0, 1, 4)
	assert.Equal(t, map[interface{}][]float64{0: {0}, 1: {1, 4}}, s.GroupBy(fn))
}

func TestFloatStreamIterate(t *testing.T) {
	fn := func(element float64) float64 {
		return element * 2
	}
	s := FloatStreamOf().Iterate(1, fn)
	next := s.First()
	assert.Equal(t, 2.0, next.MustGet())
	next = s.First()
	assert.Equal(t, 4.0, next.MustGet())
	next = s.First()
	assert.Equal(t, 8.0, next.MustGet())
}

func TestFloatStreamLast(t *testing.T) {
	s := FloatStreamOf()
	last := s.Last()
	assert.True(t, last.IsEmpty())

	s = FloatStreamOf(1)
	last = s.Last()
	assert.Equal(t, 1.0, last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())

	s = FloatStreamOf(1, 2)
	last = s.Last()
	assert.Equal(t, 2.0, last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())
}

func TestFloatStreamMap(t *testing.T) {
	fn := func(element float64) float64 {
		return element * 2
	}
	s := FloatStreamOf().Map(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = FloatStreamOf(1).Map(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = FloatStreamOf(1, 2).Map(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestFloatStreamMapToInt(t *testing.T) {
	fn := func(element float64) int {
		return int(element * 2)
	}
	s := FloatStreamOf().MapToInt(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = FloatStreamOf(1).MapToInt(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = FloatStreamOf(1, 2).MapToInt(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestFloatStreamMapTo(t *testing.T) {
	fn := func(element float64) interface{} {
		return strconv.FormatFloat(element*2, 'f', -1, 64)
	}
	s := FloatStreamOf().MapTo(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = FloatStreamOf(1).MapTo(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = FloatStreamOf(1, 2).MapTo(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestFloatStreamMapToString(t *testing.T) {
	fn := func(element float64) string {
		return strconv.FormatFloat(element*2, 'f', -1, 64)
	}
	s := FloatStreamOf().MapToString(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = FloatStreamOf(1).MapToString(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = FloatStreamOf(1, 2).MapToString(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestFloatStreamMax(t *testing.T) {
	s := FloatStreamOf()
	max := s.Max()
	assert.True(t, max.IsEmpty())

	s = FloatStreamOf(1)
	max = s.Max()
	assert.Equal(t, 1.0, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = FloatStreamOf(1, 2)
	max = s.Max()
	assert.Equal(t, 2.0, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = FloatStreamOf(1, 3, 2)
	max = s.Max()
	assert.Equal(t, 3.0, max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())
}

func TestFloatStreamMin(t *testing.T) {
	s := FloatStreamOf()
	min := s.Min()
	assert.True(t, min.IsEmpty())

	s = FloatStreamOf(1)
	min = s.Min()
	assert.Equal(t, 1.0, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = FloatStreamOf(1, 0)
	min = s.Min()
	assert.Equal(t, 0.0, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = FloatStreamOf(1, -1, 2)
	min = s.Min()
	assert.Equal(t, -1.0, min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())
}

func TestFloatStreamNoneMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := FloatStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = FloatStreamOf(3, 4)
	assert.True(t, s.NoneMatch(fn))

	s = FloatStreamOf(1, 2, 3)
	assert.False(t, s.NoneMatch(fn))

	s = FloatStreamOf(1, 2, 3, 4)
	assert.False(t, s.NoneMatch(fn))
}

func TestFloatStreamPeek(t *testing.T) {
	var elements []float64
	fn := func(element float64) {
		elements = append(elements, element)
	}
	s := FloatStreamOf().Peek(fn)
	assert.Equal(t, elements, []float64(nil), s.ToSlice())

	elements = nil
	s = FloatStreamOf(1).Peek(fn)
	assert.Equal(t, elements, []float64{1}, s.ToSlice())

	elements = nil
	s = FloatStreamOf(1, 2).Peek(fn)
	assert.Equal(t, elements, []float64{1, 2}, s.ToSlice())
}

func TestFloatStreamReduce(t *testing.T) {
	fn := func(accumulator interface{}, element float64) interface{} {
		return accumulator.(float64) + element
	}
	s := FloatStreamOf()
	sum := s.Reduce(0.0, fn)
	assert.Equal(t, 0.0, sum)

	s = FloatStreamOf(1.0, 2.0, 3.0)
	sum = s.Reduce(1.0, fn)
	assert.Equal(t, 7.0, sum)
}

func TestFloatStreamReverseSorted(t *testing.T) {
	s := FloatStreamOf().ReverseSorted()
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = FloatStreamOf(1).ReverseSorted()
	assert.Equal(t, []float64{1}, s.ToSlice())

	s = FloatStreamOf(1, 2).ReverseSorted()
	assert.Equal(t, []float64{2, 1}, s.ToSlice())

	s = FloatStreamOf(2, 3, 1).ReverseSorted()
	assert.Equal(t, []float64{3, 2, 1}, s.ToSlice())
}

func TestFloatStreamSkip(t *testing.T) {
	s := FloatStreamOf().Skip(0)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = FloatStreamOf(1).Skip(0)
	assert.Equal(t, []float64{1}, s.ToSlice())

	s = FloatStreamOf(1).Skip(1)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = FloatStreamOf(1, 2).Skip(1)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = FloatStreamOf(1, 2, 3).Skip(2)
	assert.Equal(t, []float64{3}, s.ToSlice())

	s = FloatStreamOf(1, 2, 3, 4).Skip(2)
	assert.Equal(t, []float64{3, 4}, s.ToSlice())
}

func TestFloatStreamSum(t *testing.T) {
	s := FloatStreamOf()
	sum := s.Sum()
	assert.True(t, sum.IsEmpty())

	s = FloatStreamOf(3, 4)
	sum = s.Sum()
	assert.Equal(t, 7.0, sum.MustGet())
	sum = s.Sum()
	assert.True(t, sum.IsEmpty())

	s = FloatStreamOf(1, 2, 3)
	sum = s.Sum()
	assert.Equal(t, 6.0, sum.MustGet())
	sum = s.Sum()
	assert.True(t, sum.IsEmpty())
}

func TestFloatStreamToMap(t *testing.T) {
	fn := func(element float64) (k interface{}, v interface{}) {
		return element, strconv.FormatFloat(element, 'f', -1, 64)
	}
	s := FloatStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = FloatStreamOf(1)
	assert.Equal(t, map[interface{}]interface{}{1.0: "1"}, s.ToMap(fn))

	s = FloatStreamOf(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1.0: "1", 2.0: "2", 3.0: "3"}, s.ToMap(fn))
}
