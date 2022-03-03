package d24

import (
	"math"
)

type MultiplyExpression struct {
	binaryExpression
}

type multiplyRange struct {
	lhs          Range
	rhs          Range
	cachedValues *[]int
}

func NewMultiplyExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'*',
		NewMultiplyExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &MultiplyExpression{
		binaryExpression: b,
	}
}

func (e *MultiplyExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *MultiplyExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			return lhs * rhs, nil
		},
	)
}

func (e *MultiplyExpression) Range() Range {
	return buildBinaryExpressionRange(
		"MultiplyExpression",
		&e.binaryExpression,
		func(lhs, rhs int) (int, error) {
			return lhs * rhs, nil
		},
		func(lhs int, rhs ContinuousRange) (Range, error) {
			if lhs == 0 {
				return newContinuousRange(0, 0, 1), nil
			} else if lhs == 1 {
				return rhs, nil
			} else {
				return newContinuousRange(
					rhs.Min()*lhs,
					rhs.Max()*lhs,
					rhs.Step()*int(math.Abs(float64(lhs))),
				), nil
			}
		},
		nil,
		nil,
	)

}

func (e *MultiplyExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs Expression, rhs Expression) Expression {
			lhsLiteral, lhsIsLiteral := lhs.(*LiteralExpression)
			rhsLiteral, rhsIsLiteral := rhs.(*LiteralExpression)

			if lhsIsLiteral && rhsIsLiteral {
				return NewLiteralExpression(lhsLiteral.value * rhsLiteral.value)
			} else if lhsIsLiteral && lhsLiteral.value == 0 {
				return NewLiteralExpression(0)
			} else if lhsIsLiteral && lhsLiteral.value == 1 {
				return rhs
			} else if rhsIsLiteral && rhsLiteral.value == 0 {
				return NewLiteralExpression(0)
			} else if rhsIsLiteral && rhsLiteral.value == 1 {
				return lhs
			} else {
				return NewMultiplyExpression(lhs, rhs)
			}
		},
	)
}

// Given a set of expressions being multiplied together, recursively unrolls them into
// up to 1 literal value, a list of referenced inputs, and up to 1 other expression
func unrollMultiplyExpressions(expressions ...Expression) (*LiteralExpression, []*InputExpression, Expression) {
	result := struct {
		literal *LiteralExpression
		inputs  []*InputExpression
		other   Expression
	}{}

	for _, expr := range expressions {
		switch e := expr.(type) {
		case *LiteralExpression:
			result.literal = multiplyLiterals(result.literal, e)
		case *InputExpression:
			result.inputs = append(result.inputs, e)
		case *MultiplyExpression:
			l, i, o := unrollMultiplyExpressions(e.Lhs(), e.Rhs())
			result.literal = multiplyLiterals(result.literal, l)
			result.inputs = append(result.inputs, i...)
			if result.other == nil {
				result.other = o
			} else {
				result.other = NewMultiplyExpression(result.other, o)
			}
		default:
			if result.other == nil {
				result.other = expr
			} else {
				result.other = NewMultiplyExpression(result.other, expr)
			}
		}
	}

	return result.literal, result.inputs, result.other
}
