package d24

import (
	"fmt"
)

type EqualsExpression struct {
	binaryExpression
}

func NewEqualsExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'=',
		NewEqualsExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &EqualsExpression{
		binaryExpression: b,
	}
}

func (e *EqualsExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *EqualsExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			if lhs == rhs {
				return 1, nil
			} else {
				return 0, nil
			}
		},
	)
}

func (e *EqualsExpression) Range() Range {
	return newContinuousRange(0, 1, 1)
}

func (e *EqualsExpression) Simplify(inputs map[int]int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs, rhs Expression) Expression {
			lhsRange := lhs.Range()
			rhsRange := rhs.Range()

			context := fmt.Sprintf("simplify EqualsExpression: %s", e)

			// If the ranges of each side of the comparison will never intersect,
			// then we can always return "0" for this expression
			if !RangesIntersect(lhsRange, rhsRange, context) {
				return zeroLiteral
			}

			if RangesAreEqual(lhsRange, rhsRange, context) {
				return oneLiteral
			}

			return NewEqualsExpression(lhs, rhs)
		},
	)
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}
