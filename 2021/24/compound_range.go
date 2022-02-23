package d24

import (
	"fmt"
	"sort"
	"strings"
)

type compoundRange struct {
	members []Range
}

func newCompoundRange(members ...Range) Range {

	for {
		for i := 0; i < len(members); i++ {
			if members[i] == nil {
				continue
			}

			for j := i + 1; j < len(members); j++ {
				if members[j] == nil {
					continue
				}
				combined := tryCombineRanges(members[i], members[j])
				if combined != nil {
					members[i] = combined
					members[j] = nil
				}
			}
		}

		nextMembers := make([]Range, 0, len(members))
		for _, r := range members {
			if r != nil {
				nextMembers = append(nextMembers, r)
			}
		}

		didChange := len(nextMembers) != len(members)
		members = nextMembers

		if !didChange {
			break
		}
	}

	if len(members) == 1 {
		return members[0]
	}

	sort.Slice(members, func(i, j int) bool {
		a := members[i]
		b := members[j]

		aContinuous, aIsContinuous := a.(*continuousRange)
		bContinuous, bIsContinuous := b.(*continuousRange)

		if aIsContinuous && bIsContinuous {
			return aContinuous.min < bContinuous.min
		}

		return i < j
	})

	return &compoundRange{members}
}

// Returns whether this range includes the given value
func (r *compoundRange) Includes(value int) bool {
	for _, r := range r.members {
		if r.Includes(value) {
			return true
		}
	}
	return false
}

func (r *compoundRange) Split(around Range) (Range, Range, Range) {
	panic("compoundRange.Split() not implemented")
}

func (r *compoundRange) String() string {
	items := make([]string, len(r.members))
	for i, r := range r.members {
		items[i] = r.String()
	}
	return fmt.Sprintf("<%s>", strings.Join(items, ","))
}

func (r *compoundRange) Values() func() (int, bool) {

	var next func() (int, bool)

	ranges := make([]Range, len(r.members))
	copy(ranges, r.members)

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

func tryCombineRanges(a, b Range) Range {
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
	}

	return nil
}
