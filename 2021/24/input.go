package d24

import (
	"fmt"
	"log"
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

func (e *InputExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	if target < MinInputValue || target > MaxInputValue {
		return nil, fmt.Errorf("InputExpression can only seek targets in the range %d-%d (got %d)", MinInputValue, MaxInputValue, target)
	}

	m := make(map[int]int, 1)
	m[e.index] = target
	return m, nil
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

func (e *InputExpression) String() string {
	return fmt.Sprintf("i%d", e.index)
}
