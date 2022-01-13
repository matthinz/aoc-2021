package d24

import "fmt"

type InputExpression struct {
	index int
}

const MinInputValue = 1

const MaxInputValue = 9

func NewInputExpression(index int) Expression {
	return &InputExpression{index}
}

func (e *InputExpression) Evaluate(inputs []int) int {
	return inputs[e.index]
}

func (e *InputExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	if target < MinInputValue || target > MaxInputValue {
		return nil, fmt.Errorf("InputExpression can only seek targets in the range %d-%d (got %d)", MinInputValue, MaxInputValue, target)
	}

	m := make(map[int]int, 1)
	m[e.index] = target
	return m, nil
}

func (e *InputExpression) Range() IntRange {
	return NewIntRange(MinInputValue, MaxInputValue)
}

func (e *InputExpression) Simplify() Expression {
	return e
}

func (e *InputExpression) String() string {
	return fmt.Sprintf("i%d", e.index)
}
