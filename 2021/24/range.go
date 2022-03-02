package d24

import (
	"fmt"
	"sort"
	"time"
)

type Range interface {
	// Returns whether this range includes the given value
	Includes(value int) bool

	// Splits this range into its component ContinuousRange elements.
	Split() []ContinuousRange

	String() string

	// Returns a function that, when executed, returns the next value in the
	// range along with a boolean indicating whether the call succeeded.
	Values(context string) func() (int, bool)
}

// Builds a single Range from a list of integer values
func newRangeFromInts(values interface{}) Range {

	var sortedValues []int

	switch x := values.(type) {
	case map[int]bool:
		sortedValues = make([]int, 0, len(x))
		for value := range x {
			sortedValues = append(sortedValues, value)
		}
	case map[int]int:
		sortedValues = make([]int, 0, len(x))
		for value := range x {
			sortedValues = append(sortedValues, value)
		}

	case []int:
		sortedValues = make([]int, 0, len(x))
		for _, value := range x {
			sortedValues = append(sortedValues, value)
		}

	default:
		panic(fmt.Sprintf("Invalid value passed to newRangeFromInts: %T", values))
	}

	if len(sortedValues) == 0 {
		return EmptyRange
	}

	sort.Ints(sortedValues)

	step := 0
	start := sortedValues[0]
	prev := start

	var ranges []ContinuousRange

	for i := 1; i < len(sortedValues); i++ {
		if sortedValues[i] == prev {
			continue
		}

		if step == 0 {
			step = sortedValues[i] - prev
		} else if sortedValues[i] != prev+step {
			// This is a new range
			ranges = append(ranges, newContinuousRange(start, prev, step))
			start = sortedValues[i]
			step = 0
		}

		prev = sortedValues[i]
	}

	if step == 0 {
		// have an incomplete range
		ranges = append(ranges, newContinuousRange(start, prev, (start-prev)+1))
	} else if prev != start {
		ranges = append(ranges, newContinuousRange(start, prev, step))
	}

	if len(ranges) == 1 {
		return ranges[0]
	} else {
		return NewContinuousRangeSet(ranges)
	}
}

// Reads all values in the given range into a slice and returns it.
func GetAllValuesOfRange(r Range, context string) []int {
	result := make([]int, 0)

	nextValue := r.Values(context)

	for value, ok := nextValue(); ok; value, ok = nextValue() {
		result = append(result, value)
	}
	return result
}

// If r specifies a single value, it is returned along with `true`.
func GetSingleValueOfRange(r Range) (int, bool) {

	c, isContinuous := r.(*continuousRange)
	if isContinuous && c.min == c.max {
		return c.min, true
	}

	return 0, false
}

func RangesAreEqual(a, b Range, context string) bool {
	aContinuous, aIsContinuous := a.(*continuousRange)
	bContinuous, bIsContinuous := b.(*continuousRange)

	if aIsContinuous && bIsContinuous {
		return *aContinuous == *bContinuous
	}

	const SeenInA = 1
	const SeenInB = 2

	nextA := a.Values(context)
	nextB := b.Values(context)

	seen := make(map[int]uint8)

	for {
		aValue, aOk := nextA()
		bValue, bOk := nextB()

		if !(aOk || bOk) {
			break
		}

		if aOk {
			bits := seen[aValue]
			if bits&SeenInA == 0 {
				seen[aValue] = bits | SeenInA
			}
		}

		if bOk {
			bits := seen[bValue]
			if bits&SeenInB == 0 {
				seen[bValue] = bits | SeenInB
			}
		}
	}

	for _, bits := range seen {
		if bits != SeenInA|SeenInB {
			return false
		}
	}

	return true
}

// Returns true if a and b have _any_ values in common
func RangesIntersect(a, b Range, context string) bool {

	aContinuous, aIsContinuous := a.(*continuousRange)
	bContinuous, bIsContinuous := b.(*continuousRange)

	if aIsContinuous && bIsContinuous {
		return aContinuous.Intersects(bContinuous)
	} else if aIsContinuous {
		// `Includes()` is a very efficient call for continuousRanges, so use
		// a for that
		temp := a
		a = b
		b = temp
	}

	nextA := a.Values(context)
	for aValue, ok := nextA(); ok; aValue, ok = nextA() {
		if b.Includes(aValue) {
			return true
		}
	}

	return false
}

func buildBinaryExpressionRangeValues(
	lhs Range,
	rhs Range,
	op func(lhsValue, rhsValue int) int,
	context string,
) *[]int {
	started := time.Now()

	values := make(map[int]int)
	nextLhs := lhs.Values(context)

	steps := 0

	for lhsValue, ok := nextLhs(); ok; lhsValue, ok = nextLhs() {
		nextRhs := rhs.Values(context)
		for rhsValue, ok := nextRhs(); ok; rhsValue, ok = nextRhs() {
			value := op(lhsValue, rhsValue)
			values[value]++
			steps++
		}
	}

	uniqueValues := make([]int, 0, len(values))
	for value := range values {
		uniqueValues = append(uniqueValues, value)
	}

	duration := time.Now().Sub(started)
	if duration.Seconds() > 1 {
		// fmt.Printf("buildBinaryExpressionRangeValues: %v (%d steps, %d unique values)\n", duration, steps, len(uniqueValues))
		// fmt.Printf("\t%s\n", context)
		// fmt.Printf("\tlhs: %s\n", lhs.String())
		// fmt.Printf("\trhs: %s\n", rhs.String())
	}

	return &uniqueValues
}
