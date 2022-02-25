package d24

import (
	"fmt"
	"time"
)

type Range interface {
	// Returns whether this range includes the given value
	Includes(value int) bool

	String() string

	// Returns a function that, when executed, returns the next value in the
	// range along with a boolean indicating whether the call succeeded.
	Values(context string) func() (int, bool)
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

// Returns true if two ranges contain the exact same elements.
func RangesAreEqual(a, b Range, context string) bool {

	nextA := a.Values(context)
	nextB := b.Values(context)

	for {
		aValue, aOk := nextA()
		bValue, bOk := nextB()

		if !aOk && !bOk {
			return true
		} else if !aOk || !bOk {
			return false
		}

		if aValue != bValue {
			return false
		}
	}
}

func RangesEqual(a, b Range, context string) bool {

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

		if !(aOk && bOk) {
			break
		}

		bits := seen[aValue]
		if bits&SeenInA == 0 {
			seen[aValue] = bits | SeenInA
		}

		bits = seen[bValue]
		if bits&SeenInB == 0 {
			seen[bValue] = bits | SeenInB
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
	nextA := a.Values(context)
	for aValue, ok := nextA(); ok; aValue, ok = nextA() {
		nextB := b.Values(context)
		for bValue, ok := nextB(); ok; bValue, ok = nextB() {
			if aValue == bValue {
				return true
			}
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
		fmt.Printf("buildBinaryExpressionRangeValues: %v (%d steps, %d unique values)\n", duration, steps, len(uniqueValues))
		fmt.Printf("\t%s\n", context)
		fmt.Printf("\tlhs: %s\n", lhs.String())
		fmt.Printf("\trhs: %s\n", rhs.String())
	}

	return &uniqueValues
}
