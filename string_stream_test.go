package stream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSliceIterator(t *testing.T) {
	ai := stringSliceIterator{array: []string{"1", "2", "3"}}
	next, hasNext := ai.next()
	assert.Equal(t, "1", next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, "2", next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.Equal(t, "3", next)
	assert.True(t, hasNext)
	next, hasNext = ai.next()
	assert.False(t, hasNext)
}

func TestStringStreamNewStringStream(t *testing.T) {
	ai := stringSliceIterator{array: []string{"1", "2", "3"}}
	s := NewStringStream(ai.next)
	assert.Equal(t, []string{"1", "2", "3"}, s.ToSlice())

	s = NewStringStreamOf("3", "2", "1")
	assert.Equal(t, []string{"3", "2", "1"}, s.ToSlice())
}

func TestStringStreamAllMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := NewStringStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = NewStringStreamOf("1", "2")
	assert.True(t, s.AllMatch(fn))

	s = NewStringStreamOf("1", "2", "3")
	assert.False(t, s.AllMatch(fn))

	s = NewStringStreamOf("1", "2", "3", "4")
	assert.False(t, s.AllMatch(fn))
}

func TestStringStreamAnyMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := NewStringStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = NewStringStreamOf("3", "4")
	assert.False(t, s.AnyMatch(fn))

	s = NewStringStreamOf("1", "2", "3")
	assert.True(t, s.AnyMatch(fn))
}

func TestStringStreamConcat(t *testing.T) {
	s1 := NewStringStreamOf("1", "2", "3")
	s2 := NewStringStreamOf("4", "5", "6")
	s3 := s1.Concat(s2)
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, s3.ToSlice())
}

func TestStringStreamCount(t *testing.T) {
	s := NewStringStreamOf()
	assert.Equal(t, 0, s.Count())

	s = NewStringStreamOf("2")
	assert.Equal(t, 1, s.Count())

	s = NewStringStreamOf("2", "3")
	assert.Equal(t, 2, s.Count())
}

func TestStringStreamDistinct(t *testing.T) {
	s := NewStringStreamOf()
	assert.Equal(t, []string(nil), s.Distinct().ToSlice())

	s = NewStringStreamOf("1", "1")
	assert.Equal(t, []string{"1"}, s.Distinct().ToSlice())

	s = NewStringStreamOf("1", "2", "2", "1")
	assert.Equal(t, []string{"1", "2"}, s.Distinct().ToSlice())
}

func TestStringStreamDuplicate(t *testing.T) {
	s := NewStringStreamOf()
	assert.Equal(t, []string(nil), s.Duplicate().ToSlice())

	s = NewStringStreamOf("1", "1", "2")
	assert.Equal(t, []string{"1"}, s.Duplicate().ToSlice())

	s = NewStringStreamOf("1", "2", "2", "1", "3")
	assert.Equal(t, []string{"2", "1"}, s.Duplicate().ToSlice())
}

func TestStringStreamFilter(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := NewStringStreamOf()
	assert.Equal(t, []string(nil), s.Filter(fn).ToSlice())

	s = NewStringStreamOf("1", "2", "3")
	assert.Equal(t, []string{"1", "2"}, s.Filter(fn).ToSlice())
}

func TestStringStreamFirst(t *testing.T) {
	s := NewStringStreamOf()

	s = NewStringStreamOf("1")
	next, hasNext := s.First()
	assert.Equal(t, "1", next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)

	s = NewStringStreamOf("1", "2")
	next, hasNext = s.First()
	assert.Equal(t, "1", next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.Equal(t, "2", next)
	assert.True(t, hasNext)
	next, hasNext = s.First()
	assert.False(t, hasNext)
}

func TestStringStreamForEach(t *testing.T) {
	var elements []string
	fn := func(element string) {
		elements = append(elements, element)
	}
	s := NewStringStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []string(nil), elements)

	elements = nil
	s = NewStringStreamOf("1")
	s.ForEach(fn)
	assert.Equal(t, []string{"1"}, elements)

	elements = nil
	s = NewStringStreamOf("1", "2", "3")
	s.ForEach(fn)
	assert.Equal(t, []string{"1", "2", "3"}, elements)
}

func TestStringStreamGroupBy(t *testing.T) {
	fn := func(element string) (key interface{}) {
		i, _ := strconv.Atoi(element)
		return i % 3
	}
	s := NewStringStreamOf()
	assert.Equal(t, map[interface{}][]string{}, s.GroupBy(fn))

	s = NewStringStreamOf("0")
	assert.Equal(t, map[interface{}][]string{0: []string{"0"}}, s.GroupBy(fn))

	s = NewStringStreamOf("0", "1", "4")
	assert.Equal(t, map[interface{}][]string{0: []string{"0"}, 1: []string{"1", "4"}}, s.GroupBy(fn))
}

func TestStringStreamIterate(t *testing.T) {
	fn := func(element string) string {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := NewStringStreamOf().Iterate("1", fn)
	element, _ := s.First()
	assert.Equal(t, "2", element)
	element, _ = s.First()
	assert.Equal(t, "4", element)
	element, _ = s.First()
	assert.Equal(t, "8", element)
}

func TestStringStreamLast(t *testing.T) {
	s := NewStringStreamOf()
	next, hasNext := s.Last()
	assert.False(t, hasNext)

	s = NewStringStreamOf("1")
	next, hasNext = s.Last()
	assert.Equal(t, "1", next)
	assert.True(t, hasNext)

	s = NewStringStreamOf("1", "2")
	next, hasNext = s.Last()
	assert.Equal(t, "2", next)
	assert.True(t, hasNext)
}

func TestStringStreamMap(t *testing.T) {
	fn := func(element string) string {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := NewStringStreamOf().Map(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewStringStreamOf("1").Map(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = NewStringStreamOf("1", "2").Map(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestStringStreamMapToFloat64(t *testing.T) {
	fn := func(element string) float64 {
		i, _ := strconv.ParseFloat(element, 64)
		return i * 2
	}
	s := NewStringStreamOf().MapToFloat(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = NewStringStreamOf("1").MapToFloat(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = NewStringStreamOf("1", "2").MapToFloat(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestStringStreamMapToInt(t *testing.T) {
	fn := func(element string) int {
		i, _ := strconv.Atoi(element)
		return i * 2
	}
	s := NewStringStreamOf().MapToInt(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = NewStringStreamOf("1").MapToInt(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = NewStringStreamOf("1", "2").MapToInt(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestStringStreamMapToObject(t *testing.T) {
	fn := func(element string) interface{} {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := NewStringStreamOf().MapToObject(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = NewStringStreamOf("1").MapToObject(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = NewStringStreamOf("1", "2").MapToObject(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestStringStreamMax(t *testing.T) {
	s := NewStringStreamOf()
	_, valid := s.Max()
	assert.False(t, valid)

	s = NewStringStreamOf("1")
	max, valid := s.Max()
	assert.Equal(t, "1", max)
	assert.True(t, valid)

	s = NewStringStreamOf("1", "2")
	max, valid = s.Max()
	assert.Equal(t, "2", max)
	assert.True(t, valid)

	s = NewStringStreamOf("1", "3", "2")
	max, valid = s.Max()
	assert.Equal(t, "3", max)
	assert.True(t, valid)
}

func TestStringStreamMin(t *testing.T) {
	s := NewStringStreamOf()
	_, valid := s.Min()
	assert.False(t, valid)

	s = NewStringStreamOf("1")
	min, valid := s.Min()
	assert.Equal(t, "1", min)
	assert.True(t, valid)

	s = NewStringStreamOf("1", "0")
	min, valid = s.Min()
	assert.Equal(t, "0", min)
	assert.True(t, valid)

	s = NewStringStreamOf("3", "1", "2")
	min, valid = s.Min()
	assert.Equal(t, "1", min)
	assert.True(t, valid)
}

func TestStringStreamNoneMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := NewStringStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = NewStringStreamOf("3", "4")
	assert.True(t, s.NoneMatch(fn))

	s = NewStringStreamOf("1", "2", "3")
	assert.False(t, s.NoneMatch(fn))

	s = NewStringStreamOf("1", "2", "3", "4")
	assert.False(t, s.NoneMatch(fn))
}

func TestStringStreamPeek(t *testing.T) {
	var elements []string
	fn := func(element string) {
		elements = append(elements, element)
	}
	s := NewStringStreamOf().Peek(fn)
	assert.Equal(t, elements, []string(nil), s.ToSlice())

	elements = nil
	s = NewStringStreamOf("1").Peek(fn)
	assert.Equal(t, elements, []string{"1"}, s.ToSlice())

	elements = nil
	s = NewStringStreamOf("1", "2").Peek(fn)
	assert.Equal(t, elements, []string{"1", "2"}, s.ToSlice())
}

func TestStringStreamReduce(t *testing.T) {
	fn := func(accumulator interface{}, element string) interface{} {
		return accumulator.(string) + element
	}
	s := NewStringStreamOf()
	sum := s.Reduce("0", fn)
	assert.Equal(t, "0", sum)

	s = NewStringStreamOf("1", "2", "3")
	sum = s.Reduce("1", fn)
	assert.Equal(t, "1123", sum)
}

func TestStringStreamSkip(t *testing.T) {
	s := NewStringStreamOf().Skip(0)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewStringStreamOf("1").Skip(0)
	assert.Equal(t, []string{"1"}, s.ToSlice())

	s = NewStringStreamOf("1").Skip(1)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewStringStreamOf("1", "2").Skip(1)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = NewStringStreamOf("1", "2", "3").Skip(2)
	assert.Equal(t, []string{"3"}, s.ToSlice())

	s = NewStringStreamOf("1", "2", "3", "4").Skip(2)
	assert.Equal(t, []string{"3", "4"}, s.ToSlice())
}

func TestStringStreamSorted(t *testing.T) {
	s := NewStringStreamOf().Sorted()
	assert.Equal(t, []string(nil), s.ToSlice())

	s = NewStringStreamOf("1").Sorted()
	assert.Equal(t, []string{"1"}, s.ToSlice())

	s = NewStringStreamOf("2", "1").Sorted()
	assert.Equal(t, []string{"1", "2"}, s.ToSlice())

	s = NewStringStreamOf("2", "3", "1").Sorted()
	assert.Equal(t, []string{"1", "2", "3"}, s.ToSlice())
}

func TestStringStreamToMap(t *testing.T) {
	fn := func(element string) (k interface{}, v interface{}) {
		i, _ := strconv.Atoi(element)
		return element, i
	}
	s := NewStringStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = NewStringStreamOf("1")
	assert.Equal(t, map[interface{}]interface{}{"1": 1}, s.ToMap(fn))

	s = NewStringStreamOf("1", "2", "3")
	assert.Equal(t, map[interface{}]interface{}{"1": 1, "2": 2, "3": 3}, s.ToMap(fn))
}
