package d24

import (
	"fmt"
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

	if e.cachedRange != nil {
		return e.cachedRange
	}

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	if lhsIsContinuous && rhsIsContinuous {
		if rhsContinuous.min == rhsContinuous.max {
			if lhsContinuous.step == rhsContinuous.min {
				// When the step of the lhs range == value of rhs range, then
				// then the lhs range is reduced to a single value.
				value := lhsContinuous.min % rhsContinuous.min
				e.cachedRange = newContinuousRange(value, value, 1)
				return e.cachedRange
			}
		}
	}

	uniqueValues := make(map[int]int)

	nextRhsValue := e.Rhs().Range().Values("ModuloExpression.Range")
	for rhsValue, ok := nextRhsValue(); ok; rhsValue, ok = nextRhsValue() {
		values := make(map[int]int)
		nextLhsValue := e.Lhs().Range().Values("ModuloExpression.Range")

		for lhsValue, ok := nextLhsValue(); ok; lhsValue, ok = nextLhsValue() {
			value := lhsValue % rhsValue
			values[value]++
			if values[0] >= 2 {
				break
			}
		}

		for value := range values {
			uniqueValues[value]++
		}
	}

	values := make([]int, 0, len(uniqueValues))
	for value := range uniqueValues {
		values = append(values, value)
	}

	e.cachedRange = newRangeFromInts(values)

	return e.cachedRange
}

func (e *ModuloExpression) Simplify(inputs map[int]int) Expression {
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
