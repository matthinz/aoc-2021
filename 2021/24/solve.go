package d24

import (
	"fmt"
	"log"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func SolveForLargest(expr Expression, target int, l *log.Logger) ([]int, error) {
	initialInputs := []int{}

	inputs, err := solveStep(expr, target, initialInputs, 0, countInputs(expr), MaxInputValue, MinInputValue, -1, l)

	if err != nil {
		return []int{}, err
	}

	simplified := expr.Simplify(inputs)
	value, err := simplified.Evaluate()
	if err != nil {
		return []int{}, err
	}

	if value != target {
		return inputs, fmt.Errorf("Inputs did not evaluate to %d with initial expression (got %d)", target, value)
	}

	return inputs, nil
}

func SolveForSmallest(expr Expression, target int, l *log.Logger) ([]int, error) {
	initialInputs := []int{2}

	inputs, err := solveStep(expr, target, initialInputs, 0, countInputs(expr), MinInputValue, MaxInputValue, 1, l)

	if err != nil {
		return []int{}, err
	}

	simplified := expr.Simplify(inputs)
	value, err := simplified.Evaluate()
	if err != nil {
		return []int{}, err
	}

	if value != target {
		return inputs, fmt.Errorf("Inputs did not evaluate to %d with initial expression (got %d)", target, value)
	}

	return inputs, nil
}

func countInputs(expr Expression) int {
	inputCounts := make(map[int]int)
	expr.Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputCounts[input.index]++
		}
	})
	return len(inputCounts)
}

func solveStep(expr Expression, target int, inputs []int, index int, inputCount int, inputStartValue, inputEndValue, inputStep int, l *log.Logger) ([]int, error) {

	if len(inputs) >= inputCount {
		return inputs, nil
	}

	var nextInputs []int
	if index >= len(inputs) {
		nextInputs = make([]int, index+1)
	} else {
		nextInputs = make([]int, len(inputs))
	}
	copy(nextInputs, inputs)

	var initialValue int
	if IsValidInputValue(nextInputs[index]) {
		initialValue = nextInputs[index]
	} else {
		initialValue = inputStartValue
	}

	for i := initialValue; IsValidInputValue(i); i += inputStep {
		nextInputs[index] = i

		values := nextInputs[0 : index+1]

		beforeCount := countExpressions(expr)
		simplified := expr.Simplify(values)
		afterCount := countExpressions(simplified)

		l.Printf("trying %v (simplified %d%% from %d to %d nodes)", values, int((float64(beforeCount-afterCount)/float64(beforeCount))*-100), beforeCount, afterCount)

		r := simplified.Range()

		if r.Includes(target) {

			value, err := simplified.Evaluate()

			if err == nil && value == target {
				return values, nil
			}

			result, err := solveStep(simplified, target, nextInputs, index+1, inputCount, inputStartValue, inputEndValue, inputStep, l)

			if err == nil {
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not solve for %d at %d", target, index)
}

func countExpressions(expr Expression) int {
	count := 0
	expr.Accept(func(e Expression) {
		count++
	})
	return count
}
