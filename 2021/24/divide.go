package d24

import (
	"fmt"
)

type DivideExpression struct {
	binaryExpression
}

type divisionRange struct {
	lhs, rhs     Range
	cachedValues *[]int
}

func NewDivideExpression(expressions ...interface{}) Expression {
	b := newBinaryExpression(
		'/',
		NewDivideExpression,
		expressions,
	)
	if b.rhs == nil {
		return b.lhs
	}
	return &DivideExpression{
		binaryExpression: b,
	}
}

func (e *DivideExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *DivideExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Can't divide by 0")
			}
			return lhs / rhs, nil
		},
	)
}

func (e *DivideExpression) Range() Range {
	if e.cachedRange == nil {

		lhsRange := e.lhs.Range()
		rhsRange := e.rhs.Range()

		lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
		rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

		if lhsIsContinuous && rhsIsContinuous {

			if rhsContinuous.min == rhsContinuous.max {

				if rhsContinuous.min == 0 {
					e.cachedRange = EmptyRange
				} else {

					rhsIsFactorOfLhsStep := (lhsContinuous.step/rhsContinuous.min)*rhsContinuous.min == lhsContinuous.step

					if rhsIsFactorOfLhsStep {
						// If lhs is continuous and rhs is a factor of the step of lhs,
						// then we can cleanly divide everything
						e.cachedRange = newContinuousRange(
							lhsContinuous.min/rhsContinuous.min,
							lhsContinuous.max/rhsContinuous.max,
							lhsContinuous.step/rhsContinuous.min,
						)
					}
				}
			}

		}

		if e.cachedRange == nil {

			e.cachedRange = &divisionRange{
				lhs: e.Lhs().Range(),
				rhs: e.Rhs().Range(),
			}
		}
	}

	return e.cachedRange
}

func (e *DivideExpression) Simplify(inputs map[int]int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs, rhs Expression) Expression {

			return NewDivideExpression(lhs, rhs)
		},
	)
}

func recursiveDivide(dividend Expression, divisor Expression) Expression {
	if sum, isSum := dividend.(*AddExpression); isSum {
		result := divideSum(sum, divisor)
		if result != nil {
			return result
		}
	} else if product, isProduct := dividend.(*MultiplyExpression); isProduct {
		result := divideMultiplyExpression(product, divisor)
		if result != nil {
			return result
		}
	} else if literal, isLiteral := dividend.(*LiteralExpression); isLiteral {
		result := divideLiteralExpression(literal, divisor)
		if result != nil {
			return result
		}
	} else if input, isInput := dividend.(*InputExpression); isInput {
		result := divideInputExpression(input, divisor)
		if result != nil {
			return result
		}
	}

	return NewDivideExpression(dividend, divisor)
}

func divideInputExpression(dividend *InputExpression, divisor Expression) Expression {
	switch d := divisor.(type) {
	case *InputExpression:
		if d.index == dividend.index {
			return NewLiteralExpression(1)
		}
	}
	return nil
}

func divideLiteralExpression(dividend *LiteralExpression, divisor Expression) Expression {
	switch d := divisor.(type) {
	case *LiteralExpression:
		value := dividend.value / d.value
		if value*d.value == dividend.value {
			return NewLiteralExpression(value)
		}
	}
	return nil
}

func divideMultiplyExpression(dividend *MultiplyExpression, divisor Expression) Expression {
	literal, inputs, other := unrollMultiplyExpressions(dividend)

	switch d := divisor.(type) {
	case *LiteralExpression:
		if literal != nil {
			value := literal.value / d.value
			if value*d.value == literal.value {
				return NewMultiplyExpression(inputs, other, NewLiteralExpression(value))
			}
		}
	case *InputExpression:
		found := false
		for i := range inputs {
			if inputs[i].index == d.index {
				// this one cancels
				inputs[i] = nil
				found = true
				break
			}
		}
		if found {
			return NewMultiplyExpression(literal, inputs, other)
		}
	}

	return nil
}

func divideSum(dividend *AddExpression, divisor Expression) Expression {
	switch d := divisor.(type) {
	case *AddExpression:
		if *dividend == *d {
			return NewLiteralExpression(1)
		}
	}
	return NewDivideExpression(dividend, divisor)
}

////////////////////////////////////////////////////////////////////////////////
// divisionRange

func (r *divisionRange) Includes(value int) bool {
	next := r.Values(fmt.Sprintf("%s includes %v", r, value))

	for v, ok := next(); ok; v, ok = next() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *divisionRange) String() string {
	return fmt.Sprintf("<%s / %s>", r.lhs.String(), r.rhs.String())
}

func (r *divisionRange) Values(context string) func() (int, bool) {

	pos := 0

	return func() (int, bool) {

		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue / rhsValue },
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
