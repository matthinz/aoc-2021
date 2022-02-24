package d24

import (
	"fmt"
	"io/ioutil"
	"log"
)

const ActuallyLog = true

// Take a binary expression (e.g. a +, *, /, etc.) and find inputs required
// to get it to equal <target>
func findInputsForBinaryExpression(
	b BinaryExpression,
	target int,
	getRhsValues func(lhsValue int, rhsRange Range) (chan int, error),
	d InputDecider,
	l *log.Logger,
) (map[int]int, error) {

	if !ActuallyLog {
		l = log.New(ioutil.Discard, "", log.Flags())
	}

	e := b.(Expression)

	r := e.Range()
	if !r.Includes(target) {
		return nil, fmt.Errorf("Range of expression %s (%s) does not include %d", e.String(), r.String(), target)
	}

	lhsRange := b.Lhs().Range()
	rhsRange := b.Rhs().Range()

	var best map[int]int

	// for each value in left side's range, look for a corresponding value in the
	// right side's range and figure out the inputs needed to get them both to go there
	nextLhsValue := lhsRange.Values()
	for lhsValue, ok := nextLhsValue(); ok; lhsValue, ok = nextLhsValue() {
		potentialRhsValues, err := getRhsValues(lhsValue, rhsRange)

		if err != nil {
			l.Print(err.Error())
			continue
		}

		for rhsValue := range potentialRhsValues {
			lhsInputs, err := b.Lhs().FindInputs(lhsValue, d, l)

			if err != nil {
				l.Print(err.Error())
				continue
			}

			rhsInputs, err := b.Rhs().FindInputs(rhsValue, d, l)

			if err != nil {
				l.Print(err.Error())
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
	}

	if best == nil {
		return nil, fmt.Errorf("No inputs can reach %d for %s (searching %s : %s)", target, e.String(), lhsRange.String(), rhsRange.String())
	}

	if len(best) > 4 {
		l.Printf("%s @ %d: %v", e.String(), target, best)
	}

	return best, nil
}
