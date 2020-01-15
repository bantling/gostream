package stream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamSliceIterator(t *testing.T) {
	ai := sliceIterator{array: []interface{}{1, 2, 3}}
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

func TestStreamNewStream(t *testing.T) {
	ai := sliceIterator{array: []interface{}{1, 2, 3}}
	s := NewStream(ai.next)
	assert.Equal(t, []interface{}{1, 2, 3}, s.ToSlice())

	s = NewStreamOf(3, 2, 1)
	assert.Equal(t, []interface{}{3, 2, 1}, s.ToSlice())
}

func TestStreamAllMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := NewStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = NewStreamOf(1, 2)
	assert.True(t, s.AllMatch(fn))

	s = NewStreamOf(1, 2, 3)
	assert.False(t, s.AllMatch(fn))

	s = NewStreamOf(1, 2, 3, 4)
	assert.False(t, s.AllMatch(fn))
}

func TestStreamAnyMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := NewStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = NewStreamOf(3, 4)
	assert.False(t, s.AnyMatch(fn))

	s = NewStreamOf(1, 2, 3)
	assert.True(t, s.AnyMatch(fn))
}

func TestStreamConcat(t *testing.T) {
	s1 := NewStreamOf(1, 2, 3)
	s2 := NewStreamOf(4, 5, 6)
	s3 := s1.Concat(s2)
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5, 6}, s3.ToSlice())
}

func TestStreamCount(t *testing.T) {
	s := NewStreamOf()
	assert.Equal(t, 0, s.Count())

	s = NewStreamOf(2)
	assert.Equal(t, 1, s.Count())

	s = NewStreamOf(2, 3)
	assert.Equal(t, 2, s.Count())
}

func TestStreamDistinct(t *testing.T) {
	s := NewStreamOf()
	assert.Equal(t, []interface{}(nil), s.Distinct().ToSlice())

	s = NewStreamOf(1, 1)
	assert.Equal(t, []interface{}{1}, s.Distinct().ToSlice())

	s = NewStreamOf(1, 2, 2, 1)
	assert.Equal(t, []interface{}{1, 2}, s.Distinct().ToSlice())
}

func TestStreamDuplicate(t *testing.T) {
	s := NewStreamOf()
	assert.Equal(t, []interface{}(nil), s.Duplicate().ToSlice())

	s = NewStreamOf(1, 1, 2)
	assert.Equal(t, []interface{}{1}, s.Duplicate().ToSlice())

	s = NewStreamOf(1, 2, 2, 1, 3)
	assert.Equal(t, []interface{}{2, 1}, s.Duplicate().ToSlice())
}

func TestStreamFilter(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := NewStreamOf()
	assert.Equal(t, []interface{}(nil), s.Filter(fn).ToSlice())

	s = NewStreamOf(1, 2, 3)
	assert.Equal(t, []interface{}{1, 2}, s.Filter(fn).ToSlice())
}

func TestStreamFirst(t *testing.T) {
	s := NewStreamOf()

	s = NewStreamOf(1)
	next, hasNext := s.First()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)

	s = NewStreamOf(1, 2)
	next, hasNext = s.First()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.Equal(t, 2, next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)
}

func TestStreamForEach(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := NewStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []interface{}(nil), elements)

	elements = nil
	s = NewStreamOf(1)
	s.ForEach(fn)
	assert.Equal(t, []interface{}{1}, elements)

	elements = nil
	s = NewStreamOf(1, 2, 3)
	s.ForEach(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, elements)
}

func TestStreamGroupBy(t *testing.T) {
	fn := func(element interface{}) (key interface{}) {
		return element.(int) % 3
	}
	s := NewStreamOf()
	assert.Equal(t, map[interface{}][]interface{}{}, s.GroupBy(fn))

	s = NewStreamOf(0)
	assert.Equal(t, map[interface{}][]interface{}{0: []interface{}{0}}, s.GroupBy(fn))

	s = NewStreamOf(0, 1, 4)
	assert.Equal(t, map[interface{}][]interface{}{0: []interface{}{0}, 1: []interface{}{1, 4}}, s.GroupBy(fn))
}

func TestStreamIterate(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return element.(int) * 2
	}
	s := NewStreamOf().Iterate(1, fn)
	element, _ := s.First()
	assert.Equal(t, 2, element)
	element, _ = s.First()
	assert.Equal(t, 4, element)
	element, _ = s.First()
	assert.Equal(t, 8, element)
}

func TestStreamLast(t *testing.T) {
	s := NewStreamOf()
	next, hasNext := s.Last()
	assert.False(t, hasNext)

	s = NewStreamOf(1)
	next, hasNext = s.Last()
	assert.Equal(t, 1, next)
	assert.True(t, hasNext)

	s = NewStreamOf(1, 2)
	next, hasNext = s.Last()
	assert.Equal(t, 2, next)
	assert.True(t, hasNext)
}

func TestStreamMap(t *testing.T) {
	fn := func(element interface{}) interface{} {
		return strconv.Itoa(element.(int) * 2)
	}
	s := NewStreamOf().Map(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStreamOf(1).Map(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = NewStreamOf(1, 2).Map(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestStreamMapToInt(t *testing.T) {
	fn := func(element interface{}) int {
		return element.(int) * 2
	}
	s := NewStreamOf().MapToInt(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewStreamOf(1).MapToInt(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = NewStreamOf(1, 2).MapToInt(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestStreamMapToFloat(t *testing.T) {
	fn := func(element interface{}) float64 {
		return float64(element.(int) * 2)
	}
	s := NewStreamOf().MapToFloat(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewStreamOf(1).MapToFloat(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = NewStreamOf(1, 2).MapToFloat(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestStreamMapToString(t *testing.T) {
	fn := func(element interface{}) string {
		return strconv.Itoa(element.(int) * 2)
	}
	s := NewStreamOf().MapToString(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewStreamOf(1).MapToString(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = NewStreamOf(1, 2).MapToString(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestStreamMax(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := NewStreamOf()
	_, valid := s.Max(fn)
	assert.False(t, valid)

	s = NewStreamOf(1)
	max, valid := s.Max(fn)
	assert.Equal(t, 1, max)
	assert.True(t, valid)

	s = NewStreamOf(1, 2)
	max, valid = s.Max(fn)
	assert.Equal(t, 2, max)
	assert.True(t, valid)

	s = NewStreamOf(1, 3, 2)
	max, valid = s.Max(fn)
	assert.Equal(t, 3, max)
	assert.True(t, valid)
}

func TestStreamMin(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := NewStreamOf()
	_, valid := s.Min(fn)
	assert.False(t, valid)

	s = NewStreamOf(1)
	min, valid := s.Min(fn)
	assert.Equal(t, 1, min)
	assert.True(t, valid)

	s = NewStreamOf(1, 0)
	min, valid = s.Min(fn)
	assert.Equal(t, 0, min)
	assert.True(t, valid)

	s = NewStreamOf(1, -1, 2)
	min, valid = s.Min(fn)
	assert.Equal(t, -1, min)
	assert.True(t, valid)
}

func TestStreamNoneMatch(t *testing.T) {
	fn := func(element interface{}) bool { return element.(int) < 3 }
	s := NewStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = NewStreamOf(3, 4)
	assert.True(t, s.NoneMatch(fn))

	s = NewStreamOf(1, 2, 3)
	assert.False(t, s.NoneMatch(fn))

	s = NewStreamOf(1, 2, 3, 4)
	assert.False(t, s.NoneMatch(fn))
}

func TestStreamPeek(t *testing.T) {
	var elements []interface{}
	fn := func(element interface{}) {
		elements = append(elements, element)
	}
	s := NewStreamOf().Peek(fn)
	assert.Equal(t, elements, []interface{}(nil), s.ToSlice())

	elements = nil
	s = NewStreamOf(1).Peek(fn)
	assert.Equal(t, elements, []interface{}{1}, s.ToSlice())

	elements = nil
	s = NewStreamOf(1, 2).Peek(fn)
	assert.Equal(t, elements, []interface{}{1, 2}, s.ToSlice())
}

func TestStreamReduce(t *testing.T) {
	fn := func(accumulator, element2 interface{}) interface{} {
		return accumulator.(int) + element2.(int)
	}
	s := NewStreamOf()
	sum := s.Reduce(0, fn)
	assert.Equal(t, 0, sum)

	s = NewStreamOf(1, 2, 3)
	sum = s.Reduce(1, fn)
	assert.Equal(t, 7, sum)
}

func TestStreamReverseSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := NewStreamOf().ReverseSorted(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStreamOf(1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = NewStreamOf(2, 1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{2, 1}, s.ToSlice())

	s = NewStreamOf(2, 3, 1).ReverseSorted(fn)
	assert.Equal(t, []interface{}{3, 2, 1}, s.ToSlice())
}

func TestStreamSkip(t *testing.T) {
	s := NewStreamOf().Skip(0)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStreamOf(1).Skip(0)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = NewStreamOf(1).Skip(1)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStreamOf(1, 2).Skip(1)
	assert.Equal(t, []interface{}{2}, s.ToSlice())

	s = NewStreamOf(1, 2, 3).Skip(2)
	assert.Equal(t, []interface{}{3}, s.ToSlice())

	s = NewStreamOf(1, 2, 3, 4).Skip(2)
	assert.Equal(t, []interface{}{3, 4}, s.ToSlice())
}

func TestStreamSorted(t *testing.T) {
	fn := func(element1, element2 interface{}) bool {
		return element1.(int) < element2.(int)
	}
	s := NewStreamOf().Sorted(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStreamOf(1).Sorted(fn)
	assert.Equal(t, []interface{}{1}, s.ToSlice())

	s = NewStreamOf(2, 1).Sorted(fn)
	assert.Equal(t, []interface{}{1, 2}, s.ToSlice())

	s = NewStreamOf(2, 3, 1).Sorted(fn)
	assert.Equal(t, []interface{}{1, 2, 3}, s.ToSlice())
}

func TestStreamToMap(t *testing.T) {
	fn := func(element interface{}) (k interface{}, v interface{}) {
		return element, strconv.Itoa(element.(int))
	}
	s := NewStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = NewStreamOf(1)
	assert.Equal(t, map[interface{}]interface{}{1: "1"}, s.ToMap(fn))

	s = NewStreamOf(1, 2, 3)
	assert.Equal(t, map[interface{}]interface{}{1: "1", 2: "2", 3: "3"}, s.ToMap(fn))
}
