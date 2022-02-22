package d24

type Range interface {
	// Returns whether this range includes the given value
	Includes(value int) bool

	// Splits this range around <around>. Returns 3 new Ranges representing
	// the portion of the range less than <around>, <around> itself, and the
	// portion of the range greater than <around>.
	Split(around Range) (Range, Range, Range)

	String() string

	// Returns a function that, when executed, returns the next value in the
	// range along with a boolean indicating whether the call succeeded.
	Values() func() (int, bool)
}

// Reads all values in the given range into a slice and returns it.
func GetAllValuesOfRange(r Range) []int {
	result := make([]int, 0)

	nextValue := r.Values()

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
	} else if isContinuous {
		return 0, false
	}

	nextValue := r.Values()

	var index int
	var result int

	for value, ok := nextValue(); ok; value, ok = nextValue() {
		if index > 0 {
			return 0, false
		}
		result = value
		index++
	}

	return result, true
}

// Returns true if two ranges contain the exact same elements.
func RangesAreEqual(a, b Range) bool {

	nextA := a.Values()
	nextB := b.Values()

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

func RangesEqual(a, b Range) bool {

	aContinuous, aIsContinuous := a.(*continuousRange)
	bContinuous, bIsContinuous := b.(*continuousRange)

	if aIsContinuous && bIsContinuous {
		return *aContinuous == *bContinuous
	}

	const SeenInA = 1
	const SeenInB = 2

	nextA := a.Values()
	nextB := b.Values()

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
func RangesIntersect(a, b Range) bool {
	nextA := a.Values()
	for aValue, ok := nextA(); ok; aValue, ok = nextA() {
		nextB := b.Values()
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
) *[]int {
	values := make(map[int]int)

	nextLhs := lhs.Values()

	for lhsValue, ok := nextLhs(); ok; lhsValue, ok = nextLhs() {
		nextRhs := rhs.Values()
		for rhsValue, ok := nextRhs(); ok; rhsValue, ok = nextRhs() {
			value := op(lhsValue, rhsValue)
			values[value]++
		}
	}

	uniqueValues := make([]int, 0, len(values))
	for value := range values {
		uniqueValues = append(uniqueValues, value)
	}

	return &uniqueValues
}
