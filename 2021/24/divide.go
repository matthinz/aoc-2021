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

func (e *DivideExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(dividend, divisor Expression) Expression {

			dividendRange := dividend.Range()
			if dividendRange, isContinuous := dividendRange.(*continuousRange); isContinuous {
				if dividendRange.min == 0 && dividendRange.max == 0 {
					// 0 / anything = 0
					return NewLiteralExpression(0)
				}
			}

			divisorRange := divisor.Range()
			if divisorRange, isContinuous := divisorRange.(*continuousRange); isContinuous {
				if divisorRange.min == 1 && divisorRange.max == 1 {
					// anything / 1 = anything
					return dividend
				}
			}

			var simplified Expression

			switch expr := dividend.(type) {
			case *AddExpression:
				simplified = simplifyDivisionOfAddExpression(expr, divisor, inputs)
			case *InputExpression:
				simplified = simplifyDivisionOfInputExpression(expr, divisor, inputs)
			case *LiteralExpression:
				simplified = simplifyDivisionOfLiteralExpression(expr, divisor, inputs)
			case *MultiplyExpression:
				simplified = simplifyDivisionOfMultiplyExpression(expr, divisor, inputs)
			}

			if simplified == nil {
				return NewDivideExpression(
					dividend,
					divisor,
				)
			} else {
				return simplified
			}
		},
	)
}

func simplifyDivisionOfAddExpression(dividend *AddExpression, divisor Expression, inputs []int) Expression {
	switch divisor := divisor.(type) {
	case *AddExpression:
		if *dividend == *divisor {
			return NewLiteralExpression(1)
		}
	}

	potentialNewLhs := NewDivideExpression(dividend.Lhs(), divisor)
	newLhs := potentialNewLhs.Simplify(inputs)

	potentialNewRhs := NewDivideExpression(dividend.Rhs(), divisor)
	newRhs := potentialNewRhs.Simplify(inputs)

	// If both simplifications got rid of the divide expression, consider this
	// "safe" if either one is still a divide expression, then move the divide
	// expression outside
	_, newLhsIsDivide := newLhs.(*DivideExpression)
	_, newRhsIsDivide := newRhs.(*DivideExpression)

	if newLhsIsDivide || newRhsIsDivide {
		// Abort this simplification -- we risk losing precision
		return NewDivideExpression(dividend, divisor)
	}

	return NewAddExpression(newLhs, newRhs)
}

func simplifyDivisionOfInputExpression(dividend *InputExpression, divisor Expression, inputs []int) Expression {
	switch d := divisor.(type) {
	case *LiteralExpression:
		if d.value == 1 {
			return dividend
		}
	case *InputExpression:
		if d.index == dividend.index {
			return NewLiteralExpression(1)
		}
	case *MultiplyExpression:
		// Unroll this expression and look for inputs to cancel
		literal, inputs, other := unrollMultiplyExpressions(d)
		for i := range inputs {
			if inputs[i].index == dividend.index {
				inputs[i] = nil
				return NewDivideExpression(
					NewLiteralExpression(1),
					NewMultiplyExpression(literal, inputs, other),
				)
			}
		}
	}

	return nil
}

func simplifyDivisionOfLiteralExpression(dividend *LiteralExpression, divisor Expression, inputs []int) Expression {
	if dividend.value == 0 {
		return NewLiteralExpression(0)
	}

	switch divisor := divisor.(type) {
	case *LiteralExpression:
		if divisor.value == 1 {
			return dividend
		}
		value := dividend.value / divisor.value
		if value*divisor.value == dividend.value {
			return NewLiteralExpression(value)
		}
	}
	return nil
}

func simplifyDivisionOfMultiplyExpression(dividend *MultiplyExpression, divisor Expression, inputs []int) Expression {
	literal, inputExpressions, other := unrollMultiplyExpressions(dividend)

	switch d := divisor.(type) {
	case *LiteralExpression:
		if literal != nil {
			value := literal.value / d.value
			if value*d.value == literal.value {
				return NewMultiplyExpression(inputExpressions, other, NewLiteralExpression(value)).Simplify(inputs)
			}
		}
	case *InputExpression:
		found := false
		for i := range inputExpressions {
			if inputExpressions[i].index == d.index {
				// this one cancels
				inputExpressions[i] = nil
				found = true
				break
			}
		}
		if found {
			return NewMultiplyExpression(literal, inputExpressions, other).Simplify(inputs)
		}
	}

	return nil
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
