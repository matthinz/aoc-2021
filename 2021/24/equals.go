package d24

import (
	"fmt"
)

type EqualsExpression struct {
	binaryExpression
}

func NewEqualsExpression(lhs, rhs Expression) Expression {
	return &EqualsExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '=',
		},
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

func (e *EqualsExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

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

	expr := NewEqualsExpression(lhs, rhs)
	expr.(*EqualsExpression).isSimplified = true

	return expr
}

func (e *EqualsExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewEqualsExpression(lhs, rhs)
	return expr.Simplify()
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}
