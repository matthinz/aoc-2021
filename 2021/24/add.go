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

func NewAddExpression(lhs, rhs Expression) Expression {
	return &AddExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '+',
		},
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

func (e *AddExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsLiteral, lhsInputs, lhsOther := unrollSumExpression(lhs)
	rhsLiteral, rhsInputs, rhsOther := unrollSumExpression(rhs)

	rawInputs := make([]*InputExpression, 0, len(lhsInputs)+len(rhsInputs))
	rawInputs = append(rawInputs, lhsInputs...)
	rawInputs = append(rawInputs, rhsInputs...)

	others := make([]Expression, 0)
	others = append(others, lhsOther)
	others = append(others, rhsOther)
	others = append(others, combineInputs(rawInputs...)...)

	var newLhs Expression
	for _, o := range others {
		if newLhs == nil {
			newLhs = o
		} else if o != nil {
			newLhs = NewAddExpression(newLhs, o)
		}
	}

	newRhs := sumLiterals(lhsLiteral, rhsLiteral)

	if newLhs != nil && newRhs != nil {
		expr := NewAddExpression(newLhs, newRhs)
		expr.(*AddExpression).binaryExpression.isSimplified = true
		return expr
	} else if newLhs != nil {
		return newLhs
	} else if newRhs != nil {
		return newRhs
	} else {
		return NewLiteralExpression(0)
	}
}

func (e *AddExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewAddExpression(lhs, rhs)
	return expr.Simplify()
}

// given an expression, attempts to return literal and non-literal parts that can be added together
func unrollSumExpression(expr Expression) (*LiteralExpression, []*InputExpression, Expression) {

	if literal, isLiteral := expr.(*LiteralExpression); isLiteral {
		return literal, nil, nil
	}

	if input, isInput := expr.(*InputExpression); isInput {
		return nil, []*InputExpression{input}, nil
	}

	sum, isSum := expr.(*AddExpression)
	if !isSum {
		return nil, []*InputExpression{}, expr
	}

	lhsLiteral, lhsInputs, lhsOther := unrollSumExpression(sum.Lhs())
	rhsLiteral, rhsInputs, rhsOther := unrollSumExpression(sum.Rhs())

	literal := sumLiterals(lhsLiteral, rhsLiteral)

	inputs := make([]*InputExpression, 0, len(lhsInputs)+len(rhsInputs))
	inputs = append(inputs, lhsInputs...)
	inputs = append(inputs, rhsInputs...)

	var other Expression
	if lhsOther != nil && rhsOther != nil {
		other = NewAddExpression(lhsOther, rhsOther)
	} else if lhsOther != nil {
		other = lhsOther
	} else if rhsOther != nil {
		other = rhsOther
	}

	return literal, inputs, other
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
