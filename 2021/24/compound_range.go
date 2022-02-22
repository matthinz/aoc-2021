package d24

import "fmt"

type compoundRange struct {
	a, b Range
}

func newCompoundRange(a, b Range) Range {

	aContinuous, aIsContinuous := a.(*continuousRange)
	bContinuous, bIsContinuous := b.(*continuousRange)

	if aIsContinuous && bIsContinuous && aContinuous.step == bContinuous.step {

		aInsideB := aContinuous.min >= bContinuous.min && aContinuous.max <= bContinuous.max
		if aInsideB {
			return b
		}

		bInsideA := bContinuous.min >= aContinuous.min && bContinuous.max <= aContinuous.max
		if bInsideA {
			return a
		}

		aIntersectsOrAdjacentLeft := aContinuous.min < bContinuous.min && aContinuous.max >= (bContinuous.min-bContinuous.step)
		if aIntersectsOrAdjacentLeft {
			return newContinuousRange(aContinuous.min, bContinuous.max, aContinuous.step)
		}

		bIntersectsOrAdjacentLeft := bContinuous.min < aContinuous.min && bContinuous.max >= (aContinuous.min-aContinuous.step)
		if bIntersectsOrAdjacentLeft {
			return newContinuousRange(bContinuous.min, aContinuous.max, bContinuous.step)
		}

		// Make the resulting range sorted
		if bContinuous.min < aContinuous.min {
			temp := b
			b = a
			a = temp
		}

	}

	return &compoundRange{a, b}
}

// Returns whether this range includes the given value
func (r *compoundRange) Includes(value int) bool {
	return r.a.Includes(value) || r.b.Includes(value)
}

func (r *compoundRange) Split(around Range) (Range, Range, Range) {
	panic("compoundRange.Split() not implemented")
}

func (r *compoundRange) String() string {
	return fmt.Sprintf("%s,%s", r.a.String(), r.b.String())
}

func (r *compoundRange) Values() func() (int, bool) {

	var next func() (int, bool)
	ranges := []Range{r.a, r.b}

	return func() (int, bool) {
		for {
			if next == nil {
				if len(ranges) == 0 {
					return 0, false
				}
				next = ranges[0].Values()
				ranges = ranges[1:]
			}

			value, ok := next()

			if ok {
				return value, ok
			} else {
				next = nil
			}
		}
	}
}
