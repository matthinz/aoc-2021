package d24

import (
	"fmt"
	"math"
)

type splitRelationship int

const (
	beforeSplit splitRelationship = iota
	atSplit
	afterSplit
)

// splitRange is a Range implementation that provides a subset of a parent
// range. Used to support .Split() implementations.
type splitRange struct {
	parent Range
	around Range
	rel    splitRelationship
}

func newSplitRanges(parent Range, around Range) (Range, Range, Range) {
	before := splitRange{parent, around, beforeSplit}
	at := splitRange{parent, around, atSplit}
	after := splitRange{parent, around, afterSplit}
	return &before, &at, &after
}

func (r *splitRange) Includes(value int) bool {
	// TODO: We can probably do this smarter
	nextValue := r.Values()
	for v, ok := nextValue(); ok; v, ok = nextValue() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *splitRange) Split(around Range) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *splitRange) String() string {
	switch r.rel {
	case beforeSplit:
		return fmt.Sprintf("%s less than %s", r.parent.String(), r.around.String())
	case atSplit:
		return fmt.Sprintf("intersection of %s and %s", r.parent.String(), r.around.String())
	case afterSplit:
		return fmt.Sprintf("%s greater than %s", r.parent.String(), r.around.String())
	default:
		panic("Invalid relationship")
	}
}

func (r *splitRange) Values() func() (int, bool) {
	aroundMin := math.MaxInt
	aroundMax := math.MinInt

	aroundContinuous, aroundIsContinuous := r.around.(*continuousRange)
	if aroundIsContinuous {
		aroundMin = aroundContinuous.min
		aroundMax = aroundContinuous.max
	} else {
		nextAroundValue := r.around.Values()
		for aroundValue, ok := nextAroundValue(); ok; aroundValue, ok = nextAroundValue() {
			if aroundValue < aroundMin {
				aroundMin = aroundValue
			}
			if aroundValue > aroundMax {
				aroundMax = aroundValue
			}
		}
	}

	nextValue := r.parent.Values()

	return func() (int, bool) {

		for value, ok := nextValue(); ok; value, ok = nextValue() {
			if r.rel == beforeSplit && value < aroundMin {
				return value, true
			} else if r.rel == atSplit && value >= aroundMin && value <= aroundMax {
				return value, true
			} else if r.rel == afterSplit && value > aroundMax {
				return value, true
			}
		}

		return 0, false
	}
}
