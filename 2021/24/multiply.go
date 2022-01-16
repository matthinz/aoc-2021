package d24

import "math"

type MultiplyExpression struct {
	BinaryExpression
}

func NewMultiplyExpression(lhs, rhs Expression) Expression {
	return &MultiplyExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '*',
		},
	}
}

func (e *MultiplyExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *MultiplyExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) * e.rhs.Evaluate(inputs)
}

func (e *MultiplyExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			if target == 0 {
				if lhsValue != 0 {
					// rhsValue *must* be zero
					if rhsRange.Includes(0) {
						return []int{0}, nil
					}
				}

				// lhsValue is zero, so rhsValue can be literally *any* number
				if rhsRange.Len() == 1 {
					return []int{rhsRange.min}, nil
				}

				return rhsRange.Values(), nil
			}

			if target == lhsValue {
				if rhsRange.Includes(1) {
					return []int{1}, nil
				} else {
					return []int{}, nil
				}
			}

			var result []int

			for i := rhsRange.min; i <= rhsRange.max; i++ {
				if lhsValue*i == target {
					result = append(result, i)
				}
			}

			return result, nil
		},
		d,
	)
}

func (e *MultiplyExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewIntRange(lhsRange.min*rhsRange.max, lhsRange.min*rhsRange.max)
	}

	if lhsRange.Len() == 1 {
		return NewIntRangeWithStep(
			lhsRange.min*rhsRange.min,
			lhsRange.max*rhsRange.max,
			int(math.Abs(float64(lhsRange.min))),
		)
	}

	if rhsRange.Len() == 1 {
		return NewIntRangeWithStep(
			lhsRange.min*rhsRange.min,
			lhsRange.max*rhsRange.max,
			int(math.Abs(float64(rhsRange.min))),
		)
	}

	return NewIntRange(
		lhsRange.min*rhsRange.min,
		lhsRange.max*rhsRange.max,
		// TODO: Figure out how to do step here.
	)
}

func (e *MultiplyExpression) Simplify() Expression {
	if e.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are single numbers, we are doing literal multiplication
	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewLiteralExpression(lhsRange.min * rhsRange.min)
	}

	// if either range is just "0", we'll evaluate to 0
	if (lhsRange.EqualsInt(0)) || (rhsRange.EqualsInt(0)) {
		return zeroLiteral
	}

	// if either range is just "1", we evaluate to the other
	if lhsRange.EqualsInt(1) {
		return rhs
	}

	if rhsRange.EqualsInt(1) {
		return lhs
	}

	expr := NewMultiplyExpression(lhs, rhs)
	expr.(*MultiplyExpression).isSimplified = true

	return expr
}
