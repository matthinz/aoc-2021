package d24

import (
	"fmt"
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

	if e.cachedRange != nil {
		return e.cachedRange
	}

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	if lhsIsContinuous && rhsIsContinuous {

		// When both steps are 1, we are keeping a continuous range
		if lhsContinuous.step == 1 && rhsContinuous.step == 1 {
			e.cachedRange = newContinuousRange(
				lhsContinuous.min+rhsContinuous.min,
				lhsContinuous.max+rhsContinuous.max,
				rhsContinuous.step,
			)
		}

		// If either is a single value, then that just moves the other
		if lhsContinuous.min == lhsContinuous.max {
			e.cachedRange = newContinuousRange(
				rhsContinuous.min+lhsContinuous.min,
				rhsContinuous.max+lhsContinuous.min,
				rhsContinuous.step,
			)
		} else if rhsContinuous.min == rhsContinuous.max {
			e.cachedRange = newContinuousRange(
				lhsContinuous.min+rhsContinuous.min,
				lhsContinuous.max+rhsContinuous.min,
				lhsContinuous.step,
			)
		}

		// If both are using the same step AND they're aligned (a value on either would appear on the other
		// when not considering min/max), then we can just add them directly
		if lhsContinuous.step == rhsContinuous.step {
			aligned := (rhsContinuous.min-lhsContinuous.min)%lhsContinuous.step == 0
			if aligned {
				e.cachedRange = newContinuousRange(
					lhsContinuous.min+rhsContinuous.min,
					lhsContinuous.max+rhsContinuous.max,
					lhsContinuous.step,
				)
			}
		}

	}

	if e.cachedRange == nil {
		e.cachedRange = &sumRange{
			lhs: lhsRange,
			rhs: rhsRange,
		}
	}

	return e.cachedRange
}

func (e *AddExpression) Simplify(inputs map[int]int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs, rhs Expression) Expression {
			literal, inputs, other := unrollAddExpressions(lhs, rhs)

			if literal != nil && literal.value == 0 {
				literal = nil
			}

			for _, expr := range combineInputs(inputs...) {
				if other == nil {
					other = expr
				} else {
					other = NewAddExpression(other, expr)
				}
			}

			if literal != nil && other != nil {
				return NewAddExpression(other, literal)
			} else if literal != nil {
				return literal
			} else if other != nil {
				return other
			} else {
				return NewLiteralExpression(0)
			}
		},
	)
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

////////////////////////////////////////////////////////////////////////////////
// sumRange

func (r *sumRange) Includes(value int) bool {

	lhsBounded, lhsIsBounded := r.lhs.(BoundedRange)
	rhsBounded, rhsIsBounded := r.rhs.(BoundedRange)

	if lhsIsBounded && rhsIsBounded {
		inBounds := (value >= lhsBounded.Min()+rhsBounded.Min()) && (value <= lhsBounded.Max()+rhsBounded.Max())
		if !inBounds {
			return false
		}
	}

	next := r.Values(fmt.Sprintf("%s includes %d", r, value))
	for v, ok := next(); ok; v, ok = next() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *sumRange) String() string {
	return fmt.Sprintf("<%s + %s>", r.lhs.String(), r.rhs.String())
}

func (r *sumRange) Values(context string) func() (int, bool) {

	pos := 0

	return func() (int, bool) {
		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue + rhsValue },
				context,
			)
		}

		if pos >= len(*r.cachedValues) {
			return 0, false
		}

		value := (*r.cachedValues)[pos]
		pos++

		return value, true
	}
}
