package d24

import (
	"fmt"
	"log"
	"sync"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func Solve(expr Expression, target int, l *log.Logger) (map[int]int, error) {

	type validValue struct {
		index int
		value int
	}

	inputCount := countInputs(expr)

	wg := sync.WaitGroup{}
	wg.Add(inputCount)

	ch := make(chan validValue)

	for i := 0; i < inputCount; i++ {
		go func(index int) {
			inputs := make(map[int]int, 1)
			for i := MaxInputValue; i >= MinInputValue; i-- {
				inputs[index] = i
				next := expr.Simplify(inputs)
				r := next.Range()
				if r.Includes(target) {
					ch <- validValue{index, i}
				}
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}

	panic("NOT IMPLEMENTED")

	// return solveStep(expr, target, map[int]int{}, inputCount, l)
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

func solveStep(expr Expression, target int, knownInputs map[int]int, inputCount int, l *log.Logger) (map[int]int, error) {

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

		simplified := expr.Simplify(nextInputs)
		r := simplified.Range()

		if !r.Includes(target) {
			continue
		}

		if continuous, isContinuous := r.(*continuousRange); isContinuous {
			if continuous.min == target && continuous.max == target {
				return nextInputs, nil
			}
		}

		l.Println(nextInputs)

		result, err := solveStep(simplified, target, nextInputs, inputCount, l)
		if err != nil {
			continue
		}

		return result, nil
	}

	// These inputs won't work
	return nil, fmt.Errorf("Can't solve for %d", target)
}
