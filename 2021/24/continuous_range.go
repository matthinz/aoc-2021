package d24

import "fmt"

type continuousRange struct {
	min, max int
	step     int
}

func newContinuousRange(min, max, step int) *continuousRange {
	if max < min {
		temp := max
		max = min
		min = temp
	}

	if step == 0 {
		panic("0 is not a valid step")
	}

	if (max-min)%step != 0 {
		panic(fmt.Sprintf("not possible to get to %d from %d at step %d", max, min, step))
	}

	return &continuousRange{min, max, step}
}

func (r *continuousRange) Includes(value int) bool {
	insideBounds := value >= r.min && value <= r.max
	if !insideBounds {
		return false
	}

	// But it also needs to be on the step
	isOnStep := (value-r.min)%r.step == 0
	return isOnStep
}

func (r *continuousRange) Intersects(other *continuousRange) bool {
	tests := [][2]*continuousRange{
		{r, other},
		{other, r},
	}
	for i := range tests {
		lhs, rhs := tests[i][0], tests[i][1]

		minInside := lhs.min >= rhs.min && lhs.min <= rhs.max
		if minInside {
			return (lhs.min-rhs.min)%lhs.step == 0
		}

		maxInside := lhs.max >= rhs.min && lhs.max <= rhs.max
		if maxInside {
			return (lhs.max-rhs.min)%lhs.step == 0
		}
	}

	return false
}

func (r *continuousRange) Min() int {
	return r.min
}

func (r *continuousRange) Max() int {
	return r.max
}

func (r *continuousRange) String() string {
	if r.min == r.max {
		return fmt.Sprintf("<%d>", r.min)
	} else if r.step != 1 {
		return fmt.Sprintf("<%d..%d step %d>", r.min, r.max, r.step)
	} else {
		return fmt.Sprintf("<%d..%d>", r.min, r.max)
	}
}

func (r *continuousRange) Values(context string) func() (int, bool) {
	pos := 0

	return func() (int, bool) {
		if pos > r.max-r.min {
			return 0, false
		}
		value := r.min + pos
		pos += r.step
		return value, true
	}
}
