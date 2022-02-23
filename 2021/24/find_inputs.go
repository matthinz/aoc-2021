package d24

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

const ActuallyLog = true

// Take a binary expression (e.g. a +, *, /, etc.) and find inputs required
// to get it to equal <target>
func findInputsForBinaryExpression(
	e BinaryExpression,
	target int,
	getRhsValues func(lhsValue int, rhsRange Range) (chan int, error),
	d InputDecider,
	l *log.Logger,
) (map[int]int, error) {

	if !ActuallyLog {
		l = log.New(ioutil.Discard, "", log.Flags())
	}

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	l.Printf("findInputsForBinaryExpression: %s", e.String())
	l.Printf("lhs range: %v", lhsRange.String())
	l.Printf("rhs range: %v", rhsRange.String())

	var best map[int]int

	t := reflect.TypeOf(lhsRange).Elem()
	l.Printf("lhsRange: %s", t.Name())

	// for each value in left side's range, look for a corresponding value in the
	// right side's range and figure out the inputs needed to get them both to go there
	nextLhsValue := lhsRange.Values()
	for lhsValue, ok := nextLhsValue(); ok; lhsValue, ok = nextLhsValue() {
		l.Printf("lhsValue: %d", lhsValue)

		potentialRhsValues, err := getRhsValues(lhsValue, rhsRange)

		if err != nil {
			l.Printf(err.Error())
			continue
		}

		for rhsValue := range potentialRhsValues {
			l.Printf("rhsValue: %d", rhsValue)

			lhsInputs, err := e.Lhs().FindInputs(lhsValue, d, l)

			if err != nil {
				l.Printf(err.Error())
				continue
			}

			rhsInputs, err := e.Rhs().FindInputs(rhsValue, d, l)

			if err != nil {
				l.Printf(err.Error())
				continue
			}

			l.Printf("lhsInputs: %v, rhsInputs: %v", lhsInputs, rhsInputs)

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
		return nil, fmt.Errorf("No inputs can reach %d for ranges %v and %v", target, lhsRange, rhsRange)
	}

	return best, nil
}
