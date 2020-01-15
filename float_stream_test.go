package stream

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

func TestFloatStreamNewFloatStream(t *testing.T) {
	ai := floatSliceIterator{array: []float64{1, 2, 3}}
	s := NewFloatStream(ai.next)
	assert.Equal(t, []float64{1, 2, 3}, s.ToSlice())

	s = NewFloatStreamOf(3, 2, 1)
	assert.Equal(t, []float64{3, 2, 1}, s.ToSlice())
}

func TestFloatStreamAllMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := NewFloatStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = NewFloatStreamOf(1, 2)
	assert.True(t, s.AllMatch(fn))

	s = NewFloatStreamOf(1, 2, 3)
	assert.False(t, s.AllMatch(fn))

	s = NewFloatStreamOf(1, 2, 3, 4)
	assert.False(t, s.AllMatch(fn))
}

func TestFloatStreamAnyMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := NewFloatStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = NewFloatStreamOf(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = NewFloatStreamOf(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestFloatStreamAverage(t *testing.T) {
	s := NewFloatStreamOf()
	_, haveAverage := s.Average()
	assert.False(t, haveAverage)

	s = NewFloatStreamOf(3, 4)
	average, haveAverage := s.Average()
	assert.Equal(t, 3.5, average)
	assert.True(t, haveAverage)

	s = NewFloatStreamOf(1, 2, 3)
	average, haveAverage = s.Average()
	assert.Equal(t, 2.0, average)
	assert.True(t, haveAverage)
}

func TestFloatStreamConcat(t *testing.T) {
	s1 := NewFloatStreamOf(1, 2, 3)
	s2 := NewFloatStreamOf(4, 5, 6)
	s3 := s1.Concat(s2)
	assert.Equal(t, []float64{1, 2, 3, 4, 5, 6}, s3.ToSlice())
}

func TestFloatStreamCount(t *testing.T) {
	s := NewFloatStreamOf()
	assert.Equal(t, 0, s.Count())

	s = NewFloatStreamOf(2)
	assert.Equal(t, 1, s.Count())

	s = NewFloatStreamOf(2, 3)
	assert.Equal(t, 2, s.Count())
}

func TestFloatStreamDistinct(t *testing.T) {
	s := NewFloatStreamOf()
	assert.Equal(t, []float64(nil), s.Distinct().ToSlice())

	s = NewFloatStreamOf(1, 1)
	assert.Equal(t, []float64{1}, s.Distinct().ToSlice())

	s = NewFloatStreamOf(1, 2, 2, 1)
	assert.Equal(t, []float64{1, 2}, s.Distinct().ToSlice())
}

func TestFloatStreamDuplicate(t *testing.T) {
	s := NewFloatStreamOf()
	assert.Equal(t, []float64(nil), s.Duplicate().ToSlice())

	s = NewFloatStreamOf(1, 1, 2)
	assert.Equal(t, []float64{1}, s.Duplicate().ToSlice())

	s = NewFloatStreamOf(1, 2, 2, 1, 3)
	assert.Equal(t, []float64{2, 1}, s.Duplicate().ToSlice())
}

func TestFloatStreamFilter(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := NewFloatStreamOf()
	assert.Equal(t, []float64(nil), s.Filter(fn).ToSlice())

	s = NewFloatStreamOf(1, 2, 3)
	assert.Equal(t, []float64{1, 2}, s.Filter(fn).ToSlice())
}

func TestFloatStreamFirst(t *testing.T) {
	s := NewFloatStreamOf()

	s = NewFloatStreamOf(1)
	next, hasNext := s.First()
	assert.Equal(t, 1.0, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)

	s = NewFloatStreamOf(1, 2)
	next, hasNext = s.First()
	assert.Equal(t, 1.0, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.Equal(t, 2.0, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)
}

func TestFloatStreamForEach(t *testing.T) {
	var elements []float64
	fn := func(element float64) {
		elements = append(elements, element)
	}
	s := NewFloatStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []float64(nil), elements)

	elements = nil
	s = NewFloatStreamOf(1)
	s.ForEach(fn)
	assert.Equal(t, []float64{1}, elements)

	elements = nil
	s = NewFloatStreamOf(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []float64{1, 2, 3}, elements)
}

func TestFloatStreamGroupBy(t *testing.T) {
	fn := func(element float64) (key interface{}) {
		return int(element) % 3
	}
	s := NewFloatStreamOf()
	assert.Equal(t, map[interface{}][]float64{}, s.GroupBy(fn))

	s = NewFloatStreamOf(0)
	assert.Equal(t, map[interface{}][]float64{0: []float64{0}}, s.GroupBy(fn))

	s = NewFloatStreamOf(0, 1, 4)
	assert.Equal(t, map[interface{}][]float64{0: []float64{0}, 1: []float64{1, 4}}, s.GroupBy(fn))
}

func TestFloatStreamIterate(t *testing.T) {
	fn := func(element float64) float64 {
		return element * 2
	}
	s := NewFloatStreamOf().Iterate(1, fn)
	element, _ := s.First()
	assert.Equal(t, 2.0, element)
	element, _ = s.First()
	assert.Equal(t, 4.0, element)
	element, _ = s.First()
	assert.Equal(t, 8.0, element)
}

func TestFloatStreamLast(t *testing.T) {
	s := NewFloatStreamOf()
	next, hasNext := s.Last()
	assert.False(t, hasNext)

	s = NewFloatStreamOf(1)
	next, hasNext = s.Last()
	assert.Equal(t, 1.0, next)
	assert.True(t, hasNext)

	s = NewFloatStreamOf(1, 2)
	next, hasNext = s.Last()
	assert.Equal(t, 2.0, next)
	assert.True(t, hasNext)
}

func TestFloatStreamMap(t *testing.T) {
	fn := func(element float64) float64 {
		return element * 2
	}
	s := NewFloatStreamOf().Map(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewFloatStreamOf(1).Map(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = NewFloatStreamOf(1, 2).Map(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestFloatStreamMapToInt(t *testing.T) {
	fn := func(element float64) int {
		return int(element * 2)
	}
	s := NewFloatStreamOf().MapToInt(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewFloatStreamOf(1).MapToInt(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = NewFloatStreamOf(1, 2).MapToInt(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestFloatStreamMapToObject(t *testing.T) {
	fn := func(element float64) interface{} {
		return strconv.FormatFloat(element*2, 'f', -1, 64)
	}
	s := NewFloatStreamOf().MapToObject(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewFloatStreamOf(1).MapToObject(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = NewFloatStreamOf(1, 2).MapToObject(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestFloatStreamMapToString(t *testing.T) {
	fn := func(element float64) string {
		return strconv.FormatFloat(element*2, 'f', -1, 64)
	}
	s := NewFloatStreamOf().MapToString(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewFloatStreamOf(1).MapToString(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = NewFloatStreamOf(1, 2).MapToString(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestFloatStreamMax(t *testing.T) {
	s := NewFloatStreamOf()
	_, valid := s.Max()
	assert.False(t, valid)

	s = NewFloatStreamOf(1)
	max, valid := s.Max()
	assert.Equal(t, 1.0, max)
	assert.True(t, valid)

	s = NewFloatStreamOf(1, 2)
	max, valid = s.Max()
	assert.Equal(t, 2.0, max)
	assert.True(t, valid)

	s = NewFloatStreamOf(1, 3, 2)
	max, valid = s.Max()
	assert.Equal(t, 3.0, max)
	assert.True(t, valid)
}

func TestFloatStreamMin(t *testing.T) {
	s := NewFloatStreamOf()
	_, valid := s.Min()
	assert.False(t, valid)

	s = NewFloatStreamOf(1)
	min, valid := s.Min()
	assert.Equal(t, 1.0, min)
	assert.True(t, valid)

	s = NewFloatStreamOf(1, 0)
	min, valid = s.Min()
	assert.Equal(t, 0.0, min)
	assert.True(t, valid)

	s = NewFloatStreamOf(1, -1, 2)
	min, valid = s.Min()
	assert.Equal(t, -1.0, min)
	assert.True(t, valid)
}

func TestFloatStreamNoneMatch(t *testing.T) {
	fn := func(element float64) bool { return element < 3 }
	s := NewFloatStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = NewFloatStreamOf(3, 4)
	assert.True(t, s.NoneMatch(fn))

	s = NewFloatStreamOf(1, 2, 3)
	assert.False(t, s.NoneMatch(fn))

	s = NewFloatStreamOf(1, 2, 3, 4)
	assert.False(t, s.NoneMatch(fn))
}

func TestFloatStreamPeek(t *testing.T) {
	var elements []float64
	fn := func(element float64) {
		elements = append(elements, element)
	}
	s := NewFloatStreamOf().Peek(fn)
	assert.Equal(t, elements, []float64(nil), s.ToSlice())

	elements = nil
	s = NewFloatStreamOf(1).Peek(fn)
	assert.Equal(t, elements, []float64{1}, s.ToSlice())

	elements = nil
	s = NewFloatStreamOf(1, 2).Peek(fn)
	assert.Equal(t, elements, []float64{1, 2}, s.ToSlice())
}

func TestFloatStreamReduce(t *testing.T) {
	fn := func(accumulator interface{}, element float64) interface{} {
		return accumulator.(float64) + element
	}
	s := NewFloatStreamOf()
	sum := s.Reduce(0.0, fn)
	assert.Equal(t, 0.0, sum)

	s = NewFloatStreamOf(1.0, 2.0, 3.0)
	sum = s.Reduce(1.0, fn)
	assert.Equal(t, 7.0, sum)
}

func TestFloatStreamReverseSorted(t *testing.T) {
	s := NewFloatStreamOf().ReverseSorted()
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewFloatStreamOf(1).ReverseSorted()
	assert.Equal(t, []float64{1}, s.ToSlice())

	s = NewFloatStreamOf(1, 2).ReverseSorted()
	assert.Equal(t, []float64{2, 1}, s.ToSlice())

	s = NewFloatStreamOf(2, 3, 1).ReverseSorted()
	assert.Equal(t, []float64{3, 2, 1}, s.ToSlice())
}

func TestFloatStreamSkip(t *testing.T) {
	s := NewFloatStreamOf().Skip(0)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewFloatStreamOf(1).Skip(0)
	assert.Equal(t, []float64{1}, s.ToSlice())

	s = NewFloatStreamOf(1).Skip(1)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewFloatStreamOf(1, 2).Skip(1)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = NewFloatStreamOf(1, 2, 3).Skip(2)
	assert.Equal(t, []float64{3}, s.ToSlice())

	s = NewFloatStreamOf(1, 2, 3, 4).Skip(2)
	assert.Equal(t, []float64{3, 4}, s.ToSlice())
}

func TestFloatStreamSum(t *testing.T) {
	s := NewFloatStreamOf()
	_, haveSum := s.Sum()
	assert.False(t, haveSum)

	s = NewFloatStreamOf(3, 4)
	sum, haveSum := s.Sum()
	assert.Equal(t, 7.0, sum)
	assert.True(t, haveSum)

	s = NewFloatStreamOf(1, 2, 3)
	sum, haveSum = s.Sum()
	assert.Equal(t, 6.0, sum)
	assert.True(t, haveSum)
}

func TestFloatStreamToMap(t *testing.T) {
	fn := func(element float64) (k interface{}, v interface{}) {
		return element, strconv.FormatFloat(element, 'f', -1, 64)
	}
	s := NewFloatStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = NewFloatStreamOf(1)
	assert.Equal(t, map[interface{}]interface{}{1.0: "1"}, s.ToMap(fn))

	s = NewFloatStreamOf(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1.0: "1", 2.0: "2", 3.0: "3"}, s.ToMap(fn))
}
