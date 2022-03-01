package d24

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

// Attempts to solve the given expression, returning a map of input indices to
// input values required for `expr` to evaluate to `target`.
func Solve(expr Expression, target int, l *log.Logger) (map[int]int, error) {
	return solveStep(expr, target, map[int]int{}, countInputs(expr), l)
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

	type foundValue struct {
		value int
		expr  Expression
	}

	if len(knownInputs) >= inputCount {
		return knownInputs, nil
	}

	// Spawn a goroutine for each digit and see which ones work

	wg := sync.WaitGroup{}
	wg.Add(MaxInputValue - MinInputValue + 1)

	ch := make(chan foundValue)

	for value := MaxInputValue; value >= MinInputValue; value-- {
		go func(value int) {
			nextInputs := make(map[int]int)
			maxIndex := -1
			for index, value := range knownInputs {
				nextInputs[index] = value
				if index > maxIndex {
					maxIndex = index
				}
			}

			nextInputs[maxIndex+1] = value
			simplified := expr.Simplify(nextInputs)
			r := simplified.Range()
			if r.Includes(target) {
				ch <- foundValue{
					value: value,
					expr:  simplified,
				}
			}
			wg.Done()
		}(value)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	valids := make([]foundValue, 0)
	for fv := range ch {
		valids = append(valids, fv)
	}

	sort.Slice(valids, func(i, j int) bool {
		return valids[i].value > valids[j].value
	})

	for _, fv := range valids {
		nextInputs := make(map[int]int)
		for index, v := range knownInputs {
			nextInputs[index] = v
		}

		nextInputs[len(knownInputs)] = fv.value
		l.Print(nextInputs)

		result, err := solveStep(fv.expr, target, nextInputs, inputCount, l)
		if err == nil {
			return result, nil
		}
	}

	// These inputs won't work
	return nil, fmt.Errorf("Can't solve for %d", target)
}
