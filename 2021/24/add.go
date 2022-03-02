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
			literal, inputs, other := unrollAddExpressions(lhs, rhs)

			var result Expression

			for _, expr := range combineInputs(inputs...) {
				if result == nil {
					result = expr
				} else if expr != nil {
					result = NewAddExpression(result, expr)
				}
			}

			if result == nil {
				result = other
			} else if other != nil {
				result = NewAddExpression(result, other)
			}

			if literal != nil && literal.value != 0 {
				if result == nil {
					result = literal
				} else if literal != nil {
					result = NewAddExpression(result, literal)
				}
			}

			if result == nil {
				return NewLiteralExpression(0)
			} else {
				return result
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

func tryCombineSummedInputs(expressions []Expression) []Expression {
	result := make([]Expression, 0, len(expressions))

	for i := range expressions {
		if expressions[i] == nil {
			continue
		}

		iInput, iIsInput := expressions[i].(*InputExpression)
		if !iIsInput {
			continue
		}

		multiple := 1

		for j := range expressions[i+1:] {
			if expressions[j] == nil {
				continue
			}

			jInput, jIsInput := expressions[j].(*InputExpression)
			if !jIsInput {
				continue
			}

			if iInput.index == jInput.index {
				// these two inputs can be combined
				multiple++
				expressions[j] = nil
			}
		}

		if multiple > 1 {
			result = append(result, NewMultiplyExpression(iInput, NewLiteralExpression(multiple)))
		} else {
			result = append(result, iInput)
		}
	}

	return result
}
