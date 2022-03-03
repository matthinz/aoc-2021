package d24

import (
	"fmt"
	"strconv"
)

// AddExpression defines a BinaryExpression that adds its left and righthand sides.
type AddExpression struct {
	binaryExpression
}

// sumRange is a Range implementation that represents two other Ranges
// summed together.
type sumRange struct {
	lhs, rhs     Range
	cachedValues *[]int
}

func NewAddExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'+',
		NewAddExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &AddExpression{
		binaryExpression: b,
	}
}

func (e *AddExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.Lhs().Accept(visitor)
	e.Rhs().Accept(visitor)
}

func (e *AddExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			return lhs + rhs, nil
		},
	)
}

func (e *AddExpression) Range() Range {
	return buildBinaryExpressionRange(
		"AddExpression",
		&e.binaryExpression,
		func(lhs, rhs int) (int, error) {
			return lhs + rhs, nil
		},
		func(lhs int, rhs ContinuousRange) (Range, error) {
			return newContinuousRange(
				rhs.Min()+lhs,
				rhs.Max()+lhs,
				rhs.Step(),
			), nil
		},
		nil,
		nil,
	)

}

func (e *AddExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs, rhs Expression) Expression {

			lhsLiteral, lhsIsLiteral := lhs.(*LiteralExpression)
			rhsLiteral, rhsIsLiteral := rhs.(*LiteralExpression)

			if lhsIsLiteral && rhsIsLiteral {
				return NewLiteralExpression(lhsLiteral.value + rhsLiteral.value)
			} else if lhsIsLiteral && lhsLiteral.value == 0 {
				return rhs
			} else if rhsIsLiteral && rhsLiteral.value == 0 {
				return lhs
			} else {
				return NewAddExpression(lhs, rhs)
			}
		},
	)
}

func (e *AddExpression) String() string {
	lhs := e.Lhs().String()
	rhs := e.Rhs().String()
	op := "+"

	if i, err := strconv.ParseInt(rhs, 10, 64); err == nil {
		if i < 0 {
			op = "-"
			rhs = strconv.Itoa(int(i * -1))
		}
	}

	return fmt.Sprintf("(%s %s %s)", lhs, op, rhs)
}

// Given a set of expressions being added together, recurses through them
// to find up to 1 literal value, all input references, and all other expressions.
func unrollAddExpressions(expressions ...Expression) (*LiteralExpression, []*InputExpression, Expression) {
	result := struct {
		literal *LiteralExpression
		inputs  []*InputExpression
		other   Expression
	}{}

	for _, expr := range expressions {

		if expr == nil {
			continue
		}

		if literal, isLiteral := expr.(*LiteralExpression); isLiteral {
			result.literal = sumLiterals(result.literal, literal)
			continue
		}

		if input, isInput := expr.(*InputExpression); isInput {
			result.inputs = append(result.inputs, input)
			continue
		}

		sum, isSum := expr.(*AddExpression)
		if !isSum {
			if result.other == nil {
				result.other = expr
			} else {
				result.other = NewAddExpression(result.other, expr)
			}
			continue
		}

		literal, inputs, other := unrollAddExpressions(sum.Lhs(), sum.Rhs())
		result.literal = sumLiterals(result.literal, literal)
		result.inputs = append(result.inputs, inputs...)
		if result.other == nil {
			result.other = other
		} else if other != nil {
			result.other = NewAddExpression(result.other, other)
		}
	}

	return result.literal, result.inputs, result.other
}
