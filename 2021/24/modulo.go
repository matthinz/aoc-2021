package d24

import "fmt"

type ModuloExpression struct {
	BinaryExpression
}

func NewModuloExpression(lhs, rhs Expression) Expression {
	return &ModuloExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '%',
		},
	}
}

func (e *ModuloExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) % e.rhs.Evaluate(inputs)
}

func (e *ModuloExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			// need to find rhsValues such that lhsValue % rhsValue = target
			// these would be factors of lhsValue - target
			// TODO: smarter way

			var result []int
			for i := rhsRange.min; i <= rhsRange.max; i++ {
				if lhsValue%i == target {
					result = append(result, i)
				}
			}

			return result, nil
		},
		d,
	)
}

func (e *ModuloExpression) Range() IntRange {

	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		// there is only one value
		value := lhsRange.min % rhsRange.min
		return IntRange{value, value}
	} else if lhsRange.Len() == 1 {
		// TODO
	} else if rhsRange.Len() == 1 {
		// TODO
	} else if lhsRange.LessThanRange(rhsRange) {
		// e.g. 10 % 20 = 10
		return lhsRange
	} else if rhsRange.LessThanRange(lhsRange) {
		return IntRange{
			min: rhsRange.min - 1,
			max: rhsRange.max - 1,
		}
	}

	fmt.Printf("Modulo range: %v vs %v\n", lhsRange, rhsRange)

	var min, max int

	for lhsValue := lhsRange.min; lhsValue <= lhsRange.max; lhsValue++ {
		for rhsValue := rhsRange.min; rhsValue <= rhsRange.max; rhsValue++ {

			value := lhsValue % rhsValue
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}

	return IntRange{min, max}

}

func (e *ModuloExpression) Simplify() Expression {
	if e.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if lhs is 0, we can resolve to zero
	if lhsRange.EqualsInt(0) {
		return zeroLiteral
	}

	// if both ranges are single numbers, we can evaluate to a literal
	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewLiteralExpression(lhsRange.min % rhsRange.min)
	}

	// if lhs is 1 number and *less than* the rhs range, we can eval to a literal
	if lhsRange.Len() == 1 && lhsRange.min < rhsRange.min {
		return NewLiteralExpression(lhsRange.min)
	}

	expr := NewModuloExpression(lhs, rhs)
	expr.(*ModuloExpression).isSimplified = true
	return expr
}
