package d24

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type continuousRangeSet struct {
	min, max int
	members  []ContinuousRange
}

func NewContinuousRangeSet(ranges []ContinuousRange) Range {
	members := make([]ContinuousRange, len(ranges))
	copy(members, ranges)

	sort.Slice(members, func(i, j int) bool {
		if members[i].Min() < members[j].Min() {
			return true
		} else if members[i].Min() == members[j].Min() && members[i].Max() < members[j].Max() {
			return true
		} else {
			return false
		}
	})

	optimized := optimizeSortedContinuousRangesByMerging(members)
	optimized = optimizeSortedContinuousRangesByRotating(optimized)

	min := math.MaxInt
	max := math.MinInt
	for i := range optimized {
		if optimized[i].Min() < min {
			min = optimized[i].Min()
		}
		if optimized[i].Max() > max {
			max = optimized[i].Max()
		}
	}

	return &continuousRangeSet{min, max, optimized}
}

func optimizeSortedContinuousRangesByMerging(sortedRanges []ContinuousRange) []ContinuousRange {
	const CheckInterval = 100       // every N iterations, we see if we're actually optimizing things
	const RequiredOptimization = .1 // "optimized" means the result is this % smaller than the unoptimized version

	result := make([]ContinuousRange, 0)

	for i := 0; i < len(sortedRanges); i++ {
		if i > 0 && i%CheckInterval == 0 {
			maxOptimizedLengthToContinue := i - int(float64(i)*RequiredOptimization)
			if len(result) > maxOptimizedLengthToContinue {
				// we are not actually optimizing anything
				return sortedRanges
			}
		}

		current := sortedRanges[i]

		currentMin := current.Min()
		currentMax := current.Max()

		if i == 0 {
			result = append(result, current)
			continue
		}

		for j := 0; j < len(result); j++ {
			prev := result[j]

			for prev.Includes(currentMin) && currentMin <= currentMax {
				currentMin += current.Step()
			}

			for currentMin == prev.Max()+prev.Step() && currentMin <= currentMax {
				prev = newContinuousRange(prev.Min(), currentMin, prev.Step())
				result[j] = prev
				currentMin += current.Step()
			}

			if currentMin > currentMax {
				break
			}
		}

		if currentMin == current.Min() {
			result = append(result, current)
		} else if currentMin <= currentMax {
			result = append(result, newContinuousRange(currentMin, currentMax, current.Step()))
		}
	}

	return result
}

func optimizeSortedContinuousRangesByRotating(sortedRanges []ContinuousRange) []ContinuousRange {
	// "rotating" can be used when we have a large number of ranges that share a common length, and the "step" between
	// their min and max values are in-sync

	// 1. Take our ranges and break them up by length
	rangesByLength := make(map[int][]ContinuousRange)
	for _, r := range sortedRanges {
		l := r.Length()
		rangesByLength[l] = append(rangesByLength[l], r)
	}

	result := make([]ContinuousRange, 0)

	// 2. For each group by length, further subdivided them into groups that
	//    have a common "step" between their subsequent min and max values
	for length, ranges := range rangesByLength {

		if len(ranges) <= length {
			// We won't get any benefit from processing this group
			result = append(result, ranges...)
			continue
		}

		// `ranges` is sorted by Min().
		// Break it up into chunks that move according to a fixed step
		var currentChunk []ContinuousRange
		currentStep := 0

		for i := 0; i < len(ranges); i++ {
			if len(currentChunk) == 0 {
				currentChunk = append(currentChunk, ranges[i])
				continue
			}

			prev := currentChunk[len(currentChunk)-1]
			current := ranges[i]

			stepBetweenMin := current.Min() - prev.Min()
			stepBetweenMax := current.Max() - prev.Max()

			if stepBetweenMin != stepBetweenMax {
				// This won't work at all
				result = append(result, currentChunk...)
				result = append(result, current)
				currentChunk = []ContinuousRange{}
				continue
			}

			if currentStep == 0 {
				// We're not tracking a step yet
				currentStep = stepBetweenMin
				currentChunk = append(currentChunk, current)
				continue
			}

			if stepBetweenMin == currentStep {
				currentChunk = append(currentChunk, current)
				continue
			}

			result = append(result, rotateContinuousRangeSlice(currentChunk, length, currentStep)...)
			currentChunk = []ContinuousRange{}
			currentStep = 0
		}
		result = append(result, currentChunk...)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Min() < result[j].Min() {
			return true
		} else if result[i].Min() == result[j].Min() && result[i].Max() < result[j].Max() {
			return true
		} else {
			return false
		}
	})

	return result
}

func rotateContinuousRangeSlice(slice []ContinuousRange, length int, step int) []ContinuousRange {
	if len(slice) <= length {
		// No point in rotating, as it would give us a longer slice
		return slice
	}

	result := make([]ContinuousRange, 0)

	first := slice[0]
	last := slice[len(slice)-1]

	for min := first.Min(); min <= first.Max(); min += first.Step() {
		max := last.Min() + (min - first.Min())
		result = append(
			result,
			newContinuousRange(
				min,
				max,
				step,
			),
		)
	}

	return result
}

func (c *continuousRangeSet) Includes(value int) bool {
	for _, m := range c.members {
		if m.Includes(value) {
			return true
		}
	}
	return false
}

func (c *continuousRangeSet) Min() int {
	return c.min
}

func (c *continuousRangeSet) Max() int {
	return c.max
}

func (c *continuousRangeSet) Split() []ContinuousRange {
	return c.members
}

func (c *continuousRangeSet) Values(context string) func() (int, bool) {
	var next func() (int, bool)
	index := 0

	return func() (int, bool) {
		for {
			if next == nil {
				if index >= len(c.members) {
					return 0, false
				}
				next = c.members[index].Values(context)
				index++
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

func (c *continuousRangeSet) String() string {
	if len(c.members) == 1 {
		return c.members[0].String()
	}

	items := make([]string, len(c.members))
	for i, m := range c.members {
		items[i] = m.String()
	}
	return fmt.Sprintf("<%s>", strings.Join(items, ","))
}
