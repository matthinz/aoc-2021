package d24

import (
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

func (e *LiteralExpression) Evaluate() (int, error) {
	return e.value, nil
}

func (e *LiteralExpression) Range() Range {
	return newContinuousRange(e.value, e.value, 1)
}

func (e *LiteralExpression) Simplify(inputs []int) Expression {
	return e
}

func (e *LiteralExpression) String() string {
	return strconv.FormatInt(int64(e.value), 10)
}

func multiplyLiterals(literals ...*LiteralExpression) *LiteralExpression {
	var result *LiteralExpression
	for _, l := range literals {
		if result == nil {
			result = l
		} else if l != nil {
			result = NewLiteralExpression(result.value * l.value).(*LiteralExpression)
		}
	}
	return result
}

func sumLiterals(literals ...*LiteralExpression) *LiteralExpression {
	var result *LiteralExpression
	for _, l := range literals {
		if result == nil {
			result = l
		} else if l != nil {
			result = NewLiteralExpression(l.value + result.value).(*LiteralExpression)
		}
	}
	return result
}
