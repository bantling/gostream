package stream

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
	fn := func(val int) bool { return val < 3 }
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
	fn := func(val int) bool { return val < 3 }
	s := NewIntStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = NewIntStreamOf(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestIntStreamAverage(t *testing.T) {
	s := NewIntStreamOf()
	_, haveAverage := s.Average()
	assert.False(t, haveAverage)

	s = NewIntStreamOf(3, 4)
	average, haveAverage := s.Average()
	assert.Equal(t, 3.5, average)
	assert.True(t, haveAverage)

	s = NewIntStreamOf(1, 2, 3)
	average, haveAverage = s.Average()
	assert.Equal(t, 2.0, average)
	assert.True(t, haveAverage)
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
	fn := func(val int) bool { return val < 3 }
	s := NewIntStreamOf()
	assert.Equal(t, []int(nil), s.Filter(fn).ToSlice())

	s = NewIntStreamOf(1, 2, 3)
	assert.Equal(t, []int{1, 2}, s.Filter(fn).ToSlice())
}

func TestIntStreamFirst(t *testing.T) {
	s := NewIntStreamOf()

	s = NewIntStreamOf(1)
	next, hasNext := s.First()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)

	s = NewIntStreamOf(1, 2)
	next, hasNext = s.First()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.Equal(t, 2, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)
}

func TestIntStreamForEach(t *testing.T) {
	var vals []int
	fn := func(val int) {
		vals = append(vals, val)
	}
	s := NewIntStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []int(nil), vals)

	vals = nil
	s = NewIntStreamOf(1)
	s.ForEach(fn)
	assert.Equal(t, []int{1}, vals)

	vals = nil
	s = NewIntStreamOf(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []int{1, 2, 3}, vals)
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
	fn := func(val int) int {
		return val * 2
	}
	s := NewIntStreamOf().Iterate(1, fn)
	val, _ := s.First()
	assert.Equal(t, 2, val)
	val, _ = s.First()
	assert.Equal(t, 4, val)
	val, _ = s.First()
	assert.Equal(t, 8, val)
}

func TestIntStreamLast(t *testing.T) {
	s := NewIntStreamOf()
	next, hasNext := s.Last()
	assert.False(t, hasNext)

	s = NewIntStreamOf(1)
	next, hasNext = s.Last()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)

	s = NewIntStreamOf(1, 2)
	next, hasNext = s.Last()
	assert.Equal(t, 2, next)
	assert.True(t, hasNext)
}

func TestIntStreamMap(t *testing.T) {
	fn := func(val int) interface{} {
		return strconv.Itoa(val * 2)
	}
	s := NewIntStreamOf().Map(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewIntStreamOf(1).Map(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = NewIntStreamOf(1, 2).Map(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestIntStreamMax(t *testing.T) {
	s := NewIntStreamOf()
	_, valid := s.Max()
	assert.False(t, valid)

	s = NewIntStreamOf(1)
	max, valid := s.Max()
	assert.Equal(t, 1, max)
	assert.True(t, valid)

	s = NewIntStreamOf(1, 2)
	max, valid = s.Max()
	assert.Equal(t, 2, max)
	assert.True(t, valid)

	s = NewIntStreamOf(1, 3, 2)
	max, valid = s.Max()
	assert.Equal(t, 3, max)
	assert.True(t, valid)
}

func TestIntStreamMin(t *testing.T) {
	s := NewIntStreamOf()
	_, valid := s.Min()
	assert.False(t, valid)

	s = NewIntStreamOf(1)
	min, valid := s.Min()
	assert.Equal(t, 1, min)
	assert.True(t, valid)

	s = NewIntStreamOf(1, 0)
	min, valid = s.Min()
	assert.Equal(t, 0, min)
	assert.True(t, valid)

	s = NewIntStreamOf(1, -1, 2)
	min, valid = s.Min()
	assert.Equal(t, -1, min)
	assert.True(t, valid)
}

func TestIntStreamNoneMatch(t *testing.T) {
	fn := func(val int) bool { return val < 3 }
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
	var vals []int
	fn := func(val int) {
		vals = append(vals, val)
	}
	s := NewIntStreamOf().Peek(fn)
	assert.Equal(t, vals, []int(nil), s.ToSlice())

	vals = nil
	s = NewIntStreamOf(1).Peek(fn)
	assert.Equal(t, vals, []int{1}, s.ToSlice())

	vals = nil
	s = NewIntStreamOf(1, 2).Peek(fn)
	assert.Equal(t, vals, []int{1, 2}, s.ToSlice())
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
	_, haveSum := s.Sum()
	assert.False(t, haveSum)

	s = NewIntStreamOf(3, 4)
	sum, haveSum := s.Sum()
	assert.Equal(t, 7, sum)
	assert.True(t, haveSum)

	s = NewIntStreamOf(1, 2, 3)
	sum, haveSum = s.Sum()
	assert.Equal(t, 6, sum)
	assert.True(t, haveSum)
}

func TestIntStreamToMap(t *testing.T) {
	fn := func(val int) (k interface{}, v interface{}) {
		return val, strconv.Itoa(val)
	}
	s := NewIntStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = NewIntStreamOf(1)
	assert.Equal(t, map[interface{}]interface{}{1: "1"}, s.ToMap(fn))

	s = NewIntStreamOf(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1: "1", 2: "2", 3: "3"}, s.ToMap(fn))
}
