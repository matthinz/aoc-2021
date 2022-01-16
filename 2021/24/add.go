package d24

import (
	"fmt"
	"strings"
)

// AddExpression defines a BinaryExpression that adds its left and righthand sides.
type AddExpression struct {
	BinaryExpression
}

func NewAddExpression(lhs, rhs Expression) Expression {
	return &AddExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '+',
		},
	}
}

func (e *AddExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) + e.rhs.Evaluate(inputs)
}

func (e *AddExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			rhsValue := target - lhsValue
			if rhsValue < rhsRange.min || rhsValue > rhsRange.max {
				return []int{}, nil
			}
			return []int{rhsValue}, nil
		},
		d,
	)
}

func (e *AddExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	return NewIntRange(
		lhsRange.min+rhsRange.min,
		lhsRange.max+rhsRange.max,
	)
}

func (e *AddExpression) Simplify() Expression {
	if e.BinaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are single numbers we are adding two literals
	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewLiteralExpression(lhsRange.min + rhsRange.min)
	}

	// if either range is zero, use the other
	if lhsRange.EqualsInt(0) {
		return rhs
	}

	if rhsRange.EqualsInt(0) {
		return lhs
	}

	result := NewAddExpression(lhs, rhs)
	result.(*AddExpression).isSimplified = true

	return result
}

func (e *AddExpression) String() string {
	rhsRange := e.rhs.Range()
	if rhsRange.LessThanInt(0) {
		return fmt.Sprintf("(%s - %s)", e.lhs.String(), strings.Replace(e.rhs.String(), "-", "", 1))
	} else {
		return fmt.Sprintf("(%s + %s)", e.lhs.String(), e.rhs.String())
	}
}
