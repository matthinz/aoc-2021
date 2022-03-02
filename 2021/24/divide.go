package d24

import (
	"fmt"
)

type DivideExpression struct {
	binaryExpression
}

type divisionRange struct {
	lhs, rhs     Range
	cachedValues *[]int
}

func NewDivideExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'/',
		NewDivideExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &DivideExpression{
		binaryExpression: b,
	}
}

func (e *DivideExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *DivideExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Can't divide by 0")
			}
			return lhs / rhs, nil
		},
	)
}

func (e *DivideExpression) Range() Range {
	return buildBinaryExpressionRange(
		"DivideExpression",
		&e.binaryExpression,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Can't divide by zero")
			}

			return lhs / rhs, nil
		},
		func(lhs int, rhs ContinuousRange) (Range, error) {

			if lhs == 0 {
				return newContinuousRange(0, 0, 1), nil
			}

			values := make(map[int]bool)
			nextRhs := rhs.Values("DivideExpression.Range")
			for rhs, ok := nextRhs(); ok; rhs, ok = nextRhs() {
				if rhs == 0 {
					continue
				}
				values[lhs/rhs] = true
			}

			return newRangeFromInts(values), nil
		},
		func(lhs ContinuousRange, rhs int) (Range, error) {
			if rhs == 0 {
				return nil, fmt.Errorf("Can't divide by zero")
			}

			if rhs == 1 {
				return lhs, nil
			}

			values := make(map[int]bool)
			nextLhs := lhs.Values("DivideExpression.Range")
			for lhs, ok := nextLhs(); ok; lhs, ok = nextLhs() {
				if lhs == 0 {
					continue
				}
				values[lhs/rhs] = true
			}

			return newRangeFromInts(values), nil
		},
		nil,
	)
}

func (e *DivideExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(dividend, divisor Expression) Expression {

			// We have to be _very_ careful about what we simplify here.
			// This is integer division, so many simplification rules will not apply.

			dividendRange := dividend.Range()
			if dividendRange, isContinuous := dividendRange.(*continuousRange); isContinuous {
				if dividendRange.min == 0 && dividendRange.max == 0 {
					// 0 / anything = 0
					return NewLiteralExpression(0)
				}
			}

			divisorRange := divisor.Range()
			if divisorRange, isContinuous := divisorRange.(*continuousRange); isContinuous {
				if divisorRange.min == 1 && divisorRange.max == 1 {
					// thing / 1 = thing
					return dividend
				}
			}

			// Two literals mean we can just do the division
			literalDividend, dividendIsLiteral := dividend.(*LiteralExpression)
			if dividendIsLiteral {
				literalDivisor, divisorIsLiteral := divisor.(*LiteralExpression)
				if divisorIsLiteral {
					return NewLiteralExpression(literalDividend.value / literalDivisor.value)
				}
			}

			return NewDivideExpression(dividend, divisor)
		},
	)
}

func simplifyDivisionOfLiteralExpression(dividend *LiteralExpression, divisor Expression, inputs []int) Expression {
	if dividend.value == 0 {
		return NewLiteralExpression(0)
	}

	switch divisor := divisor.(type) {
	case *LiteralExpression:
		if divisor.value == 1 {
			return dividend
		}
		value := dividend.value / divisor.value
		if value*divisor.value == dividend.value {
			return NewLiteralExpression(value)
		}
	}
	return nil
}
