package d24

import "fmt"

type EqualsExpression struct {
	BinaryExpression
}

func NewEqualsExpression(lhs, rhs Expression) Expression {
	return &EqualsExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '=',
		},
	}
}

func (e *EqualsExpression) Evaluate(inputs []int) int {
	lhsValue := e.lhs.Evaluate(inputs)
	rhsValue := e.rhs.Evaluate(inputs)
	if lhsValue == rhsValue {
		return 1
	} else {
		return 0
	}
}

func (e *EqualsExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	if target != 0 && target != 1 {
		return nil, fmt.Errorf("EqualsExpression can't seek a target other than 0 or 1 (got %d)", target)
	}

	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			if target == 0 {
				// We must find *any* rhsValue that does not equal lhsValue
				if rhsRange.EqualsInt(lhsValue) {
					return []int{}, nil
				} else if rhsRange.Len() == 1 {
					return []int{rhsRange.min}, nil
				}

				result := make([]int, 0, rhsRange.Len())
				for i := rhsRange.min; i <= rhsRange.max; i++ {
					if i != lhsValue {
						result = append(result, i)
					}
				}

				return result, nil
			}

			// We must find a value in rhsRange that equals lhsValue
			if rhsRange.Includes(lhsValue) {
				return []int{lhsValue}, nil
			}

			return []int{lhsValue}, nil
		},
		d,
	)
}

func (e *EqualsExpression) Range() IntRange {
	return NewIntRange(0, 1)
}

func (e *EqualsExpression) Simplify() Expression {
	if e.BinaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if all elements of both ranges are equal, we are comparing two equal values
	if lhsRange == rhsRange {
		return oneLiteral
	}

	// if the ranges of each side of the comparison will never intersect,
	// then we can always return "0" for this expression

	if !lhsRange.IntersectsRange(rhsRange) {
		return zeroLiteral
	}

	expr := NewEqualsExpression(lhs, rhs)
	expr.(*EqualsExpression).isSimplified = true

	return expr
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}

func (e *EqualsExpression) Visit(v func(e Expression)) {
	v(e)
	e.lhs.Visit(v)
	e.rhs.Visit(v)
}
