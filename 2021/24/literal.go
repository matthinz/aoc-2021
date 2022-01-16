package d24

import (
	"fmt"
	"log"
	"strconv"
)

type LiteralExpression struct {
	value int
}

var zeroLiteral = NewLiteralExpression(0)
var oneLiteral = NewLiteralExpression(1)

func NewLiteralExpression(value int) Expression {
	return &LiteralExpression{value}
}

func (e *LiteralExpression) Accept(visitor func(e Expression)) {
	visitor(e)
}

func (e *LiteralExpression) Evaluate(inputs []int) int {
	return e.value
}

func (e *LiteralExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	if e.value != target {
		return nil, fmt.Errorf("LiteralExpression %d can't seek target value %d", e.value, target)
	}
	// no inputs can affect this expression's value
	return map[int]int{}, nil
}

func (e *LiteralExpression) Range() IntRange {
	return NewIntRange(e.value, e.value)
}

func (e *LiteralExpression) Simplify() Expression {
	return e
}

func (e *LiteralExpression) String() string {
	return strconv.FormatInt(int64(e.value), 10)
}
