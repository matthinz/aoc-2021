package d24

import (
	"fmt"
)

type InputExpression struct {
	index int
}

const MinInputValue = 1

const MaxInputValue = 9

var inputRange = newContinuousRange(MinInputValue, MaxInputValue, 1)

func NewInputExpression(index int) Expression {
	return &InputExpression{index}
}

func (e *InputExpression) Accept(visitor func(e Expression)) {
	visitor(e)
}

func (e *InputExpression) Evaluate(inputs []int) int {
	return inputs[e.index]
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
	return fmt.Sprintf("i%d", e.index)
}
