package d24

// IntRange defines an inclusive range of integer values.
type IntRange struct {
	min, max int
}

func NewIntRange(min, max int) IntRange {
	if min > max {
		temp := max
		max = min
		min = temp
	}

	return IntRange{min, max}
}

func (r IntRange) EqualsInt(value int) bool {
	return r.min == r.max && r.min == value
}

func (r *IntRange) Includes(value int) bool {
	return value >= r.min && value <= r.max
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

func (r *IntRange) Values() []int {
	result := make([]int, r.Len())
	for i := 0; i < r.Len(); i++ {
		result[i] = r.min + i
	}
	return result
}
