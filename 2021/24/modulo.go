package d24

import (
	"fmt"
	"math"
)

// Things about modulus
//
// x % y = "remainder of x/y"
//
// This works out to: (decimal part of x / y) * y
//
// When x < y, there is *only* a decimal part, so this becomes (x/y*y) == x
//
// When x == y, then modulus == 0
//
// When x > y:
//
// 	0 <= result <= y - 1
//
// When x == 0, result is *always* 0
//
// When x < 0 and y >= 0 result is negative
// When x >= 0 and y < 0 result is positive
// When x < 0 and y < 0 result is negative

type ModuloExpression struct {
	binaryExpression
}

func NewModuloExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'%',
		NewModuloExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &ModuloExpression{
		binaryExpression: b,
	}
}

func (e *ModuloExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *ModuloExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Cannot take modulo 0")
			}
			return lhs % rhs, nil
		},
	)
}

func (e *ModuloExpression) Range() Range {

	return buildBinaryExpressionRange(
		"ModuloExpression",
		&e.binaryExpression,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Can't modulo 0")
			}
			return lhs % rhs, nil
		},
		func(lhs int, rhs ContinuousRange) (Range, error) {
			if lhs == 0 {
				return newContinuousRange(0, 0, 1), nil
			}

			values := make(map[int]int)
			nextRhs := rhs.Values("ModuloExpression.Range")

			for rhs, ok := nextRhs(); ok; rhs, ok = nextRhs() {
				if rhs == 0 {
					continue
				}
				values[lhs%rhs]++
				if values[0] > 2 {
					break
				}
			}

			return newRangeFromInts(values), nil
		},
		func(lhs ContinuousRange, rhs int) (Range, error) {
			if rhs == 0 {
				return nil, fmt.Errorf("Can't modulo 0")
			}

			maxNumberOfValues := int(math.Abs(float64(rhs)))

			values := make(map[int]bool)

			nextLhs := lhs.Values("ModuloExpression.Range")
			for lhs, ok := nextLhs(); ok; lhs, ok = nextLhs() {
				values[lhs%rhs] = true
				if len(values) >= maxNumberOfValues {
					break
				}
			}

			return newRangeFromInts(values), nil
		},
		nil,
	)
}

func (e *ModuloExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs, rhs Expression) Expression {
			lhsRange := lhs.Range()
			rhsRange := rhs.Range()

			lhsSingleValue, lhsIsSingleValue := GetSingleValueOfRange(lhsRange)
			rhsSingleValue, rhsIsSingleValue := GetSingleValueOfRange(rhsRange)

			if rhsIsSingleValue && rhsSingleValue == 0 {
				// TODO: This is invalid
				return NewModuloExpression(lhs, rhs)
			}

			// If lhs is 0, we can resolve to zero
			if lhsIsSingleValue && lhsSingleValue == 0 {
				return zeroLiteral
			}

			// If both ranges are single numbers, we can simplify to a literal
			if lhsIsSingleValue && rhsIsSingleValue {
				return NewLiteralExpression(lhsSingleValue % rhsSingleValue)
			}

			// TODO: If lhs is 1 number and *less than* the rhs range, we can eval to a literal

			return NewModuloExpression(lhs, rhs)
		},
	)
}
