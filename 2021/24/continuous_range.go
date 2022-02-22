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

	if (max-min)%step != 0 {
		panic(fmt.Sprintf("not possible to get to %d from %d at step %d", max, min, step))
	}

	return &continuousRange{min, max, step}
}

func (r *continuousRange) Includes(value int) bool {
	return value >= r.min && value <= r.max
}

func (r *continuousRange) Split(around Range) (Range, Range, Range) {

	aroundContinuous, aroundIsContinuous := around.(*continuousRange)

	if aroundIsContinuous {
		return splitContinuousWithContinuous(r, aroundContinuous)
	}

	panic("NOT IMPLEMENTED")
}

func (r *continuousRange) String() string {
	if r.min == r.max {
		return fmt.Sprintf("%d", r.min)
	} else if r.step != 1 {
		return fmt.Sprintf("%d..%d step %d", r.min, r.max, r.step)
	} else {
		return fmt.Sprintf("%d..%d", r.min, r.max)
	}
}

func (r *continuousRange) Values() func() (int, bool) {
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

func splitContinuousWithContinuous(r *continuousRange, around *continuousRange) (Range, Range, Range) {
	// We know we're splitting around a range with a defined min and max,
	// so we can take some shortcuts.

	var before, at, after Range

	// Since `r` and `around` may be at different steps, we need to run through
	// and find the parts of `around` that actually sync up with `r`'s step.

	beforeMax := around.min - 1
	for {
		if beforeMax < r.min {
			// we don't actually have a `before`, so to speak
			before = nil
			break
		}

		inR := (beforeMax-r.min)%r.step == 0
		if inR {
			before = newContinuousRange(r.min, beforeMax, r.step)
			break
		}

		beforeMax--
	}

	var atMin = around.min
	var atMax = around.max
	for {
		if atMin > atMax {
			// we don't have an `at` to speak of
			at = nil
			break
		}

		if atMin > r.max {
			at = nil
			break
		}

		if atMax < r.min {
			at = nil
			break
		}

		minInR := (atMin-r.min)%r.step == 0
		if !minInR {
			atMin++
		}

		maxInR := (atMax-r.min)%r.step == 0
		if !maxInR {
			atMax--
		}

		if minInR && maxInR {
			at = newContinuousRange(atMin, atMax, r.step)
			break
		}
	}

	afterMin := around.max + 1
	for {
		if afterMin > r.max {
			// no after
			after = nil
			break
		}

		minInR := (afterMin-r.min)%r.step == 0
		if minInR {
			after = newContinuousRange(afterMin, r.max, r.step)
			break
		}

		afterMin++
	}
	return before, at, after

}
