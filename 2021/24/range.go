package d24

type Range interface {
	// Returns whether this range includes the given value
	Includes(value int) bool

	// Splits this range around <around>. Returns 3 new Ranges representing
	// the portion of the range less than <around>, <around> itself, and the
	// portion of the range greater than <around>.
	Split(around Range) (Range, Range, Range)

	String() string

	// Returns a channel that outputs the values in this range.
	// Order of output is not guaranteed, and the output may contain duplicates.
	Values() chan int
}

// Reads all values in the given range into a slice and returns it.
func GetAllValuesOfRange(r Range) []int {
	result := make([]int, 0)
	for value := range r.Values() {
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

	values := r.Values()

	var index int
	var result int

	for value := range values {
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

	aValues := a.Values()
	bValues := b.Values()

	for {
		aValue, aOk := <-aValues
		bValue, bOk := <-bValues

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

	aValues := a.Values()
	bValues := b.Values()

	seen := make(map[int]uint8)

	for {
		aValue, aOk := <-aValues
		bValue, bOk := <-bValues

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
	aValues := a.Values()
	for aValue := range aValues {
		bValues := b.Values()
		for bValue := range bValues {
			if aValue == bValue {
				return true
			}
		}
	}
	return false
}
