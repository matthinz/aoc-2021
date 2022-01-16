package d24

import (
	"fmt"
	"math"
)

// IntRange defines an inclusive range of integer values.
type IntRange struct {
	min, max int
	step     int
}

func NewIntRange(min, max int) IntRange {
	return NewIntRangeWithStep(min, max, 1)
}

func NewIntRangeWithStep(min, max, step int) IntRange {
	if min > max {
		temp := max
		max = min
		min = temp
	}

	return IntRange{min, max, step}
}

func (r IntRange) Each(f func(i int) bool) {
	for i := r.min; i <= r.max; i += r.step {
		keepGoing := f(i)
		if !keepGoing {
			break
		}
	}
}

func (r IntRange) EqualsInt(value int) bool {
	return r.min == r.max && r.min == value
}

func (r *IntRange) Includes(value int) bool {
	return value >= r.min && value <= r.max
}

// Returns the intersection of two ranges, or an error if they do not intersect
func (r IntRange) Intersect(other IntRange) (IntRange, error) {

	if r.max < other.min || other.max < r.min {
		return IntRange{}, fmt.Errorf("Ranges do not intersect")
	}

	if r.min < other.min && r.max > other.max {
		return other, nil
	}

	if other.min < r.min && other.max > r.max {
		return r, nil
	}

	i := NewIntRange(
		int(math.Max(float64(r.min), float64(other.min))),
		int(math.Min(float64(r.max), float64(other.max))),
	)

	return i, nil
}

func (r IntRange) IntersectsRange(other IntRange) bool {
	ok := r.min >= other.min && r.min <= other.max
	ok = ok || r.max >= other.min && r.max <= other.max
	ok = ok || other.min >= r.min && other.min <= r.max
	ok = ok || other.max >= r.min && other.max <= r.max
	return ok
}

func (r IntRange) Len() int {
	return (r.max - r.min) + 1
}

func (r IntRange) LessThanInt(value int) bool {
	return r.max < value
}

func (r IntRange) LessThanRange(other IntRange) bool {
	return r.max < other.min
}

func (r IntRange) Split(around int) (*IntRange, *IntRange, *IntRange) {
	if around < r.min {
		return nil, nil, &r
	}

	if around > r.max {
		return &r, nil, nil
	}

	if around == r.min {
		if r.max == r.min {
			return nil, &r, nil
		}
		b := NewIntRange(around, around)
		c := NewIntRange(r.min+1, r.max)
		return nil, &b, &c
	}

	if around == r.max {
		if r.max == r.min {
			return nil, &r, nil
		}
		a := NewIntRange(r.min, r.max-1)
		b := NewIntRange(around, around)
		return &a, &b, nil
	}

	a := NewIntRange(r.min, around-1)
	b := NewIntRange(around, around)
	c := NewIntRange(around+1, r.max)

	return &a, &b, &c
}

func (r *IntRange) Values() []int {
	result := make([]int, r.Len())
	for i := 0; i < r.Len(); i++ {
		result[i] = r.min + i
	}
	return result
}
