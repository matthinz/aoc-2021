package d24

import (
	"fmt"
	"log"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func Solve(expr Expression, target int, l *log.Logger) ([]int, error) {
	inputs, err := solveStep(expr, target, []int{}, countInputs(expr), l)

	if err != nil {
		return []int{}, err
	}

	simplified := expr.Simplify(inputs)
	value, err := simplified.Evaluate()
	if err != nil {
		return []int{}, err
	}

	if value != target {
		return inputs, fmt.Errorf("Inputs did not evaluate to %d", target)
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

func solveStep(expr Expression, target int, inputs []int, inputCount int, l *log.Logger) ([]int, error) {

	if len(inputs) >= inputCount {
		return inputs, nil
	}

	nextInputs := make([]int, len(inputs)+1)
	copy(nextInputs, inputs)

	index := len(nextInputs) - 1

	for i := MaxInputValue; i >= MinInputValue; i-- {
		nextInputs[index] = i

		l.Printf("trying %d = %d", index, i)
		simplified := expr.Simplify(nextInputs)
		r := simplified.Range()
		l.Printf("range: %s", r)

		if r.Includes(target) {
			fmt.Println(nextInputs)
			result, err := solveStep(simplified, target, nextInputs, inputCount, l)
			if err == nil {
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not solve for %d", target)
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
