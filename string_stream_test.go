package gostream

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

func TestStringStreamStringStreamFromIter(t *testing.T) {
	ai := stringSliceIterator{array: []string{"1", "2", "3"}}
	s := StringStreamFromIter(ai.next)
	assert.Equal(t, []string{"1", "2", "3"}, s.ToSlice())

	s = StringStreamOf("3", "2", "1")
	assert.Equal(t, []string{"3", "2", "1"}, s.ToSlice())
}

func TestStringStreamAllMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := StringStreamOf()
	assert.True(t, s.AllMatch(fn))

	s = StringStreamOf("1", "2")
	assert.True(t, s.AllMatch(fn))

	s = StringStreamOf("1", "2", "3")
	assert.False(t, s.AllMatch(fn))

	s = StringStreamOf("1", "2", "3", "4")
	assert.False(t, s.AllMatch(fn))
}

func TestStringStreamAnyMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := StringStreamOf()
	assert.False(t, s.AnyMatch(fn))

	s = StringStreamOf("3", "4")
	assert.False(t, s.AnyMatch(fn))

	s = StringStreamOf("1", "2", "3")
	assert.True(t, s.AnyMatch(fn))
}

func TestStringStreamConcat(t *testing.T) {
	s1 := StringStreamOf("1", "2", "3")
	s2 := StringStreamOf("4", "5", "6")
	s3 := s1.Concat(s2)
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, s3.ToSlice())
}

func TestStringStreamCount(t *testing.T) {
	s := StringStreamOf()
	assert.Equal(t, 0, s.Count())

	s = StringStreamOf("2")
	assert.Equal(t, 1, s.Count())

	s = StringStreamOf("2", "3")
	assert.Equal(t, 2, s.Count())
}

func TestStringStreamDistinct(t *testing.T) {
	s := StringStreamOf()
	assert.Equal(t, []string(nil), s.Distinct().ToSlice())

	s = StringStreamOf("1", "1")
	assert.Equal(t, []string{"1"}, s.Distinct().ToSlice())

	s = StringStreamOf("1", "2", "2", "1")
	assert.Equal(t, []string{"1", "2"}, s.Distinct().ToSlice())
}

func TestStringStreamDuplicate(t *testing.T) {
	s := StringStreamOf()
	assert.Equal(t, []string(nil), s.Duplicate().ToSlice())

	s = StringStreamOf("1", "1", "2")
	assert.Equal(t, []string{"1"}, s.Duplicate().ToSlice())

	s = StringStreamOf("1", "2", "2", "1", "3")
	assert.Equal(t, []string{"2", "1"}, s.Duplicate().ToSlice())
}

func TestStringStreamFilter(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := StringStreamOf()
	assert.Equal(t, []string(nil), s.Filter(fn).ToSlice())

	s = StringStreamOf("1", "2", "3")
	assert.Equal(t, []string{"1", "2"}, s.Filter(fn).ToSlice())
}

func TestStringStreamFirst(t *testing.T) {
	s := StringStreamOf()
	first := s.First()
	assert.True(t, first.IsEmpty())

	s = StringStreamOf("1")
	first = s.First()
	assert.Equal(t, "1", first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())

	s = StringStreamOf("1", "2")
	first = s.First()
	assert.Equal(t, "1", first.MustGet())
	first = s.First()
	assert.Equal(t, "2", first.MustGet())
	first = s.First()
	assert.True(t, first.IsEmpty())
}

func TestStringStreamForEach(t *testing.T) {
	var elements []string
	fn := func(element string) {
		elements = append(elements, element)
	}
	s := StringStreamOf()
	s.ForEach(fn)
	assert.Equal(t, []string(nil), elements)

	elements = nil
	s = StringStreamOf("1")
	s.ForEach(fn)
	assert.Equal(t, []string{"1"}, elements)

	elements = nil
	s = StringStreamOf("1", "2", "3")
	s.ForEach(fn)
	assert.Equal(t, []string{"1", "2", "3"}, elements)
}

func TestStringStreamGroupBy(t *testing.T) {
	fn := func(element string) (key interface{}) {
		i, _ := strconv.Atoi(element)
		return i % 3
	}
	s := StringStreamOf()
	assert.Equal(t, map[interface{}][]string{}, s.GroupBy(fn))

	s = StringStreamOf("0")
	assert.Equal(t, map[interface{}][]string{0: {"0"}}, s.GroupBy(fn))

	s = StringStreamOf("0", "1", "4")
	assert.Equal(t, map[interface{}][]string{0: {"0"}, 1: {"1", "4"}}, s.GroupBy(fn))
}

func TestStringStreamIterate(t *testing.T) {
	fn := func(element string) string {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := StringStreamOf().Iterate("1", fn)
	first := s.First()
	assert.Equal(t, "2", first.MustGet())
	first = s.First()
	assert.Equal(t, "4", first.MustGet())
	first = s.First()
	assert.Equal(t, "8", first.MustGet())
}

func TestStringStreamLast(t *testing.T) {
	s := StringStreamOf()
	last := s.Last()
	assert.True(t, last.IsEmpty())

	s = StringStreamOf("1")
	last = s.Last()
	assert.Equal(t, "1", last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())

	s = StringStreamOf("1", "2")
	last = s.Last()
	assert.Equal(t, "2", last.MustGet())
	last = s.Last()
	assert.True(t, last.IsEmpty())
}

func TestStringStreamMap(t *testing.T) {
	fn := func(element string) string {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := StringStreamOf().Map(fn)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = StringStreamOf("1").Map(fn)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = StringStreamOf("1", "2").Map(fn)
	assert.Equal(t, []string{"2", "4"}, s.ToSlice())
}

func TestStringStreamMapTo(t *testing.T) {
	fn := func(element string) interface{} {
		i, _ := strconv.Atoi(element)
		return strconv.Itoa(i * 2)
	}
	s := StringStreamOf().MapTo(fn)
	assert.Equal(t, []interface{}(nil), s.ToSlice())

	s = StringStreamOf("1").MapTo(fn)
	assert.Equal(t, []interface{}{"2"}, s.ToSlice())

	s = StringStreamOf("1", "2").MapTo(fn)
	assert.Equal(t, []interface{}{"2", "4"}, s.ToSlice())
}

func TestStringStreamMapToFloat64(t *testing.T) {
	fn := func(element string) float64 {
		i, _ := strconv.ParseFloat(element, 64)
		return i * 2
	}
	s := StringStreamOf().MapToFloat(fn)
	assert.Equal(t, []float64(nil), s.ToSlice())

	s = StringStreamOf("1").MapToFloat(fn)
	assert.Equal(t, []float64{2}, s.ToSlice())

	s = StringStreamOf("1", "2").MapToFloat(fn)
	assert.Equal(t, []float64{2, 4}, s.ToSlice())
}

func TestStringStreamMapToInt(t *testing.T) {
	fn := func(element string) int {
		i, _ := strconv.Atoi(element)
		return i * 2
	}
	s := StringStreamOf().MapToInt(fn)
	assert.Equal(t, []int(nil), s.ToSlice())

	s = StringStreamOf("1").MapToInt(fn)
	assert.Equal(t, []int{2}, s.ToSlice())

	s = StringStreamOf("1", "2").MapToInt(fn)
	assert.Equal(t, []int{2, 4}, s.ToSlice())
}

func TestStringStreamMax(t *testing.T) {
	s := StringStreamOf()
	max := s.Max()
	assert.True(t, max.IsEmpty())

	s = StringStreamOf("1")
	max = s.Max()
	assert.Equal(t, "1", max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = StringStreamOf("1", "2")
	max = s.Max()
	assert.Equal(t, "2", max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())

	s = StringStreamOf("1", "3", "2")
	max = s.Max()
	assert.Equal(t, "3", max.MustGet())
	max = s.Max()
	assert.True(t, max.IsEmpty())
}

func TestStringStreamMin(t *testing.T) {
	s := StringStreamOf()
	min := s.Min()
	assert.True(t, min.IsEmpty())

	s = StringStreamOf("1")
	min = s.Min()
	assert.Equal(t, "1", min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = StringStreamOf("1", "0")
	min = s.Min()
	assert.Equal(t, "0", min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())

	s = StringStreamOf("3", "1", "2")
	min = s.Min()
	assert.Equal(t, "1", min.MustGet())
	min = s.Min()
	assert.True(t, min.IsEmpty())
}

func TestStringStreamNoneMatch(t *testing.T) {
	fn := func(element string) bool { return element < "3" }
	s := StringStreamOf()
	assert.True(t, s.NoneMatch(fn))

	s = StringStreamOf("3", "4")
	assert.True(t, s.NoneMatch(fn))

	s = StringStreamOf("1", "2", "3")
	assert.False(t, s.NoneMatch(fn))

	s = StringStreamOf("1", "2", "3", "4")
	assert.False(t, s.NoneMatch(fn))
}

func TestStringStreamPeek(t *testing.T) {
	var elements []string
	fn := func(element string) {
		elements = append(elements, element)
	}
	s := StringStreamOf().Peek(fn)
	assert.Equal(t, elements, []string(nil), s.ToSlice())

	elements = nil
	s = StringStreamOf("1").Peek(fn)
	assert.Equal(t, elements, []string{"1"}, s.ToSlice())

	elements = nil
	s = StringStreamOf("1", "2").Peek(fn)
	assert.Equal(t, elements, []string{"1", "2"}, s.ToSlice())
}

func TestStringStreamReduce(t *testing.T) {
	fn := func(accumulator interface{}, element string) interface{} {
		return accumulator.(string) + element
	}
	s := StringStreamOf()
	sum := s.Reduce("0", fn)
	assert.Equal(t, "0", sum)

	s = StringStreamOf("1", "2", "3")
	sum = s.Reduce("1", fn)
	assert.Equal(t, "1123", sum)
}

func TestStringStreamReverseSorted(t *testing.T) {
	s := StringStreamOf().ReverseSorted()
	assert.Equal(t, []string(nil), s.ToSlice())

	s = StringStreamOf("1").ReverseSorted()
	assert.Equal(t, []string{"1"}, s.ToSlice())

	s = StringStreamOf("1", "2").ReverseSorted()
	assert.Equal(t, []string{"2", "1"}, s.ToSlice())

	s = StringStreamOf("2", "3", "1").ReverseSorted()
	assert.Equal(t, []string{"3", "2", "1"}, s.ToSlice())
}

func TestStringStreamSkip(t *testing.T) {
	s := StringStreamOf().Skip(0)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = StringStreamOf("1").Skip(0)
	assert.Equal(t, []string{"1"}, s.ToSlice())

	s = StringStreamOf("1").Skip(1)
	assert.Equal(t, []string(nil), s.ToSlice())

	s = StringStreamOf("1", "2").Skip(1)
	assert.Equal(t, []string{"2"}, s.ToSlice())

	s = StringStreamOf("1", "2", "3").Skip(2)
	assert.Equal(t, []string{"3"}, s.ToSlice())

	s = StringStreamOf("1", "2", "3", "4").Skip(2)
	assert.Equal(t, []string{"3", "4"}, s.ToSlice())
}

func TestStringStreamSorted(t *testing.T) {
	s := StringStreamOf().Sorted()
	assert.Equal(t, []string(nil), s.ToSlice())

	s = StringStreamOf("1").Sorted()
	assert.Equal(t, []string{"1"}, s.ToSlice())

	s = StringStreamOf("2", "1").Sorted()
	assert.Equal(t, []string{"1", "2"}, s.ToSlice())

	s = StringStreamOf("2", "3", "1").Sorted()
	assert.Equal(t, []string{"1", "2", "3"}, s.ToSlice())
}

func TestStringStreamToMap(t *testing.T) {
	fn := func(element string) (k interface{}, v interface{}) {
		i, _ := strconv.Atoi(element)
		return element, i
	}
	s := StringStreamOf()
	assert.Equal(t, map[interface{}]interface{}{}, s.ToMap(fn))

	s = StringStreamOf("1")
	assert.Equal(t, map[interface{}]interface{}{"1": 1}, s.ToMap(fn))

	s = StringStreamOf("1", "2", "3")
	assert.Equal(t, map[interface{}]interface{}{"1": 1, "2": 2, "3": 3}, s.ToMap(fn))
}
