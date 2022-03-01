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

			// We have to be _very_ careful about what we simplify here.
			// This is integer division, so many simplification rules will not apply.

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
					// thing / 1 = thing
					return dividend
				}
			}

			// Two literals mean we can just do the division
			literalDividend, dividendIsLiteral := dividend.(*LiteralExpression)
			if dividendIsLiteral {
				literalDivisor, divisorIsLiteral := divisor.(*LiteralExpression)
				if divisorIsLiteral {
					return NewLiteralExpression(literalDividend.value / literalDivisor.value)
				}
			}

			return NewDivideExpression(dividend, divisor)
		},
	)
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
