package d24

import (
	_ "embed"
	"fmt"
	"io"
	"log"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(24, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	l.Printf("parsing...")

	reg := parseInput(r)

	l.Printf("parsed")

	inputCounts := make(map[int]int)
	reg.z.Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputCounts[input.index]++
		}
	})

	l.Printf("Found %d inputs", len(inputCounts))

	inputs, err := solve(reg.z, map[int]int{}, len(inputCounts), l)

	if err != nil {
		panic(err)
	}

	digits := make([]int, len(inputCounts))

	for i := 0; i < len(inputCounts); i++ {
		digit, isSet := inputs[i]
		if isSet {
			digits[i] = digit
		} else {
			digits[i] = 9
		}
	}

	l.Println(digits)

	result := reg.z.Evaluate(digits)

	l.Println(result)

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

func solve(expr Expression, knownInputs map[int]int, inputCount int, l *log.Logger) (map[int]int, error) {

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

		result, err := solve(simplified, nextInputs, inputCount, l)
		if err != nil {
			continue
		}

		return result, nil
	}

	// These inputs won't work
	return nil, fmt.Errorf("Can't solve expression")
}
