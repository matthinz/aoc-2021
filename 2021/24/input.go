package d24

import (
	"fmt"
	"sort"
)

type InputExpression struct {
	index int
}

const MinInputValue = 1

const MaxInputValue = 9

var inputRange = newContinuousRange(MinInputValue, MaxInputValue, 1)

var inputsByIndex = make(map[int]*InputExpression)

func NewInputExpression(index int) Expression {
	expr, ok := inputsByIndex[index]
	if ok {
		return expr
	}

	expr = &InputExpression{index}
	inputsByIndex[index] = expr

	return expr
}

func (e *InputExpression) Accept(visitor func(e Expression)) {
	visitor(e)
}

func (e *InputExpression) Evaluate() (int, error) {
	return 0, fmt.Errorf("Unknown input: %d", e.index)
}

func (e *InputExpression) Includes(value int) bool {
	return value >= MinInputValue && value <= MaxInputValue
}

func (e *InputExpression) Max() int {
	return MaxInputValue
}

func (e *InputExpression) Min() int {
	return MinInputValue
}

func (e *InputExpression) Range() Range {
	return inputRange
}

func (e *InputExpression) Simplify() Expression {
	return e
}

func (e *InputExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	value, ok := inputs[e.index]
	if !ok {
		return e
	}
	return NewLiteralExpression(value)
}

func (e *InputExpression) String() string {
	if e == nil {
		panic("e is nil")
	}
	return fmt.Sprintf("i%d", e.index)
}

// Attempts to combine inputs into either InputExpression or MultiplyExpressions
func combineInputs(inputs ...*InputExpression) []Expression {
	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].index < inputs[j].index
	})

	result := make([]Expression, 0, len(inputs))

	inputsLen := len(inputs)
	for i := 0; i < inputsLen; i++ {
		if inputs[i] == nil {
			continue
		}

		multiple := 1

		for j := i + 1; j < inputsLen; j++ {
			if inputs[j] == nil {
				continue
			}
			if inputs[j].index == inputs[i].index {
				multiple++
				inputs[j] = nil
			}
		}

		if multiple == 1 {
			result = append(result, inputs[i])
		} else {
			result = append(result, NewMultiplyExpression(inputs[i], NewLiteralExpression(multiple)))
		}
	}

	return result
}
