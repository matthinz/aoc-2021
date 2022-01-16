package d24

import (
	"fmt"
	"math"
)

type DivideExpression struct {
	BinaryExpression
}

func NewDivideExpression(lhs, rhs Expression) Expression {
	return &DivideExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '/',
		},
	}
}

func (e *DivideExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *DivideExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) / e.rhs.Evaluate(inputs)
}

func (e *DivideExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(dividend int, divisorRange IntRange) ([]int, error) {

			// 4th grade math recap: dividend / divisor = target
			// here we return potential divisors between min and max that will equal target

			if target == 0 {
				// When target == 0, divisor can't affect the result, except when it
				// can. We're doing integer division, so a large enough divisor *could* get us to zero
				// e.g if we're doing 6 / x = 0, any x > 6 will result in 0
				i, err := NewIntRange(dividend+1, math.MaxInt).Intersect(divisorRange)
				if err != nil {
					return nil, err
				}
				return i.Values(), nil
			}

			// dividend = divisor * target
			// divisor = dividend / target
			divisor := dividend / target

			if divisor == 0 {
				return nil, fmt.Errorf("Can't divide by zero")
			}

			if divisor < divisorRange.min || divisor > divisorRange.max {
				return []int{}, nil
			}

			if dividend/divisor != target {
				return []int{}, nil
			}

			return []int{divisor}, nil
		},
		d,
	)
}

func (e *DivideExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()
	return NewIntRange(
		lhsRange.min/rhsRange.max,
		lhsRange.max/rhsRange.min,
	)
}

func (e *DivideExpression) Simplify() Expression {
	if e.BinaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are single numbers, we can sub in a literal
	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewLiteralExpression(lhsRange.min / rhsRange.min)
	}

	// if left value is zero, this will eval to zero
	if lhsRange.EqualsInt(0) {
		return NewLiteralExpression(0)
	}

	// if right value is 1, this will eval to lhs
	if rhsRange.EqualsInt(1) {
		return lhs
	}

	result := NewDivideExpression(lhs, rhs)
	result.(*DivideExpression).isSimplified = true

	return result
}
