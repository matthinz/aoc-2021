package d24

import (
	"fmt"
	"log"
)

// Take a binary expression (e.g. a +, *, /, etc.) and find inputs required
// to get it to equal <target>
func findInputsForBinaryExpression(
	e BinaryExpression,
	target int,
	getRhsValues func(lhsValue int, rhsRange IntRange) ([]int, error),
	d InputDecider,
	l *log.Logger,
) (map[int]int, error) {

	l.Println()
	l.Printf("findInputsForBinaryExpression: %s", e.String())

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	l.Printf("lhsRange: %v, rhsRange: %v", lhsRange, rhsRange)

	var best map[int]int

	// for each value in left side's range, look for a corresponding value in the
	// right side's range and figure out the inputs needed to get them both to go there
	lhsRange.Each(func(lhsValue int) bool {
		potentialRhsValues, err := getRhsValues(lhsValue, rhsRange)

		if err != nil {
			return true
		}

		for _, rhsValue := range potentialRhsValues {

			lhsInputs, err := e.Lhs().FindInputs(lhsValue, d, l)

			if err != nil {
				continue
			}

			rhsInputs, err := e.Rhs().FindInputs(rhsValue, d, l)

			if err != nil {
				continue
			}

			bothSidesInSync := true
			inputs := make(map[int]int, len(lhsInputs)+len(rhsInputs))

			for index, value := range rhsInputs {
				lhsInputValue, lhsUsesInput := (lhsInputs)[index]
				if lhsUsesInput && lhsInputValue != value {
					// for this to work, left and right side need the same input set to
					// different values
					bothSidesInSync = false
					break
				}
				inputs[index] = value
			}

			if !bothSidesInSync {
				continue
			}

			for index, value := range lhsInputs {
				inputs[index] = value
			}

			if best == nil {
				best = inputs
			} else {
				b, err := d(best, inputs)
				if err == nil {
					best = b
				}
			}
		}

		return true
	})

	if best == nil {
		return nil, fmt.Errorf("No inputs can reach %d for ranges %v and %v", target, lhsRange, rhsRange)
	}

	return best, nil
}
