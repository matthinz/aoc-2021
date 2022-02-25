package d24

import (
	"fmt"
	"log"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func Solve(expr Expression, target int, l *log.Logger) (map[int]int, error) {
	inputCounts := make(map[int]int)
	expr.Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputCounts[input.index]++
		}
	})
	return nil, fmt.Errorf("Could not solve for %d", target)
}

func solveStep(expr Expression, knownInputs map[int]int, inputCount int, l *log.Logger) (map[int]int, error) {

	nextInputs := make(map[int]int)
	maxIndex := -1
	for index, value := range knownInputs {
		nextInputs[index] = value
		if index > maxIndex {
			maxIndex = index
		}
	}

	if maxIndex+1 >= inputCount {
		return knownInputs, nil
	}

	for value := MaxInputValue; value >= MinInputValue; value-- {
		nextInputs[maxIndex+1] = value

		simplified := expr.SimplifyUsingPartialInputs(nextInputs)
		r := simplified.Range()

		if !r.Includes(0) {
			continue
		}

		if continuous, isContinuous := r.(*continuousRange); isContinuous {
			if continuous.min == 0 && continuous.max == 0 {
				return nextInputs, nil
			}
		}

		l.Println(nextInputs)

		result, err := solveStep(simplified, nextInputs, inputCount, l)
		if err != nil {
			continue
		}

		return result, nil
	}

	// These inputs won't work
	return nil, fmt.Errorf("Can't solve expression")
}
