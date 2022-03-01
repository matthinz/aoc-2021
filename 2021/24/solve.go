package d24

import (
	"fmt"
	"log"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func Solve(expr Expression, target int, l *log.Logger) ([]int, error) {
	initialInputs := []int{9}

	inputs, err := solveStep(expr, target, initialInputs, 0, countInputs(expr), l)

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

	// Now we want to count _up_  to try and
	ch := solveUp(expr, target, inputs, l)
	var solution []int
	for solution = range ch {
		l.Print(solution)
	}

	return solution, nil
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

func solveStep(expr Expression, target int, inputs []int, index int, inputCount int, l *log.Logger) ([]int, error) {

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
		initialValue = MaxInputValue
	}

	for i := initialValue; i >= MinInputValue; i-- {
		nextInputs[index] = i

		values := nextInputs[0 : index+1]

		l.Printf("trying %v", values)

		beforeCount := countExpressions(expr)
		simplified := expr.Simplify(values)
		afterCount := countExpressions(simplified)

		l.Printf("Simplified from %d -> %d (%d%%)", beforeCount, afterCount, int((float64(beforeCount-afterCount)/float64(beforeCount))*-100))

		if afterCount < 50 {
			l.Printf(simplified.String())
		}

		r := simplified.Range()

		l.Printf("range is now %s", r)

		if r.Includes(target) {

			fmt.Println(values)

			value, err := simplified.Evaluate()

			if err == nil && value == target {
				return values, nil
			}

			result, err := solveStep(simplified, target, nextInputs, index+1, inputCount, l)

			if err == nil {
				return result, nil
			}
		} else {
			l.Printf("range does not include %d", target)
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

func solveUp(expr Expression, target int, inputs []int, l *log.Logger) chan []int {
	ch := make(chan []int)

	go func() {
		defer close(ch)

		index := len(inputs) - 1
		for {

			for inputs[index] >= MaxInputValue && index >= 0 {
				index--
			}

			if index < 0 {
				break
			}

			inputs[index]++

			simplified := expr.Simplify(inputs)
			result, err := simplified.Evaluate()

			if err == nil && result == target {
				inputCopy := make([]int, len(inputs))
				copy(inputCopy, inputs)
				ch <- inputCopy
			}
		}

	}()

	return ch
}
