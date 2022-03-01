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

			if lhs == rhs {
				return NewLiteralExpression(1)
			}

			lhsValue, lhsError := lhs.Evaluate()
			rhsValue, rhsError := rhs.Evaluate()

			if lhsError == nil && rhsError == nil {
				if lhsValue == rhsValue {
					return NewLiteralExpression(1)
				} else {
					return NewLiteralExpression(0)
				}
			}

			if lhs, lhsIsBounded := lhs.Range().(BoundedRange); lhsIsBounded {
				if rhs, rhsIsBounded := rhs.Range().(BoundedRange); rhsIsBounded {

					rangesCouldIntersect := ((lhs.Min() >= rhs.Min() && lhs.Min() <= rhs.Max()) ||
						(lhs.Max() >= rhs.Min() && lhs.Max() <= rhs.Max()) ||
						(rhs.Min() >= lhs.Min() && rhs.Min() <= lhs.Max()) ||
						(rhs.Max() >= lhs.Min() && rhs.Max() <= lhs.Max()))

					if !rangesCouldIntersect {
						return NewLiteralExpression(0)
					}
				}
			}

			return NewEqualsExpression(lhs, rhs)
		},
	)
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}
