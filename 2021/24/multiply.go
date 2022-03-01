package d24

import (
	"fmt"
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

	if e.cachedRange != nil {
		return e.cachedRange
	}

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	if lhsIsContinuous && rhsIsContinuous {
		if lhsContinuous.min == lhsContinuous.max {
			if lhsContinuous.min == 0 {
				e.cachedRange = newContinuousRange(0, 0, 1)
			} else {
				e.cachedRange = newContinuousRange(
					lhsContinuous.min*rhsContinuous.min,
					lhsContinuous.max*rhsContinuous.max,
					lhsContinuous.min,
				)
			}
		} else if rhsContinuous.min == rhsContinuous.max {
			if rhsContinuous.min == 0 {
				e.cachedRange = newContinuousRange(0, 0, 1)
			} else {
				e.cachedRange = newContinuousRange(
					lhsContinuous.min*rhsContinuous.min,
					lhsContinuous.max*rhsContinuous.max,
					rhsContinuous.min,
				)
			}
		} else if rhsContinuous.min == 0 && rhsContinuous.max == 1 {
			if lhsContinuous.Includes(0) {
				e.cachedRange = lhsContinuous
			} else {
				e.cachedRange = newCompoundRange(
					lhsContinuous,
					newContinuousRange(0, 0, 1),
				)
			}
		} else if lhsContinuous.min == 0 && lhsContinuous.max == 1 {
			if rhsContinuous.Includes(0) {
				e.cachedRange = rhsContinuous
			} else {
				e.cachedRange = newCompoundRange(
					rhsContinuous,
					newContinuousRange(0, 0, 1),
				)
			}
		}
	}

	if e.cachedRange == nil {
		e.cachedRange = &multiplyRange{
			lhs: lhsRange,
			rhs: rhsRange,
		}
	}

	return e.cachedRange
}

func (e *MultiplyExpression) Simplify(inputs []int) Expression {
	return simplifyBinaryExpression(
		&e.binaryExpression,
		inputs,
		func(lhs Expression, rhs Expression) Expression {
			lhsRange, rhsRange := lhs.Range(), rhs.Range()

			lhsSingleValue, lhsIsSingleValue := GetSingleValueOfRange(lhsRange)
			rhsSingleValue, rhsIsSingleValue := GetSingleValueOfRange(rhsRange)

			// if both ranges are single numbers, we are doing literal multiplication
			if lhsIsSingleValue && rhsIsSingleValue {
				return NewLiteralExpression(lhsSingleValue * rhsSingleValue)
			}

			// if either range is just "0", we'll evaluate to 0
			if lhsIsSingleValue && lhsSingleValue == 0 {
				return zeroLiteral
			}

			if rhsIsSingleValue && rhsSingleValue == 0 {
				return zeroLiteral
			}

			// if either range is just "1", we evaluate to the other
			if lhsIsSingleValue && lhsSingleValue == 1 {
				return rhs
			}

			if rhsIsSingleValue && rhsSingleValue == 1 {
				return lhs
			}

			/*
				// For expressions like (i0 + 4) * 5, distribute the "5" into  so we get
				// ((5*i0) + 20)
				lhsSum, lhsIsSum := lhs.(*AddExpression)
				rhsLiteral, rhsIsLiteral := rhs.(*LiteralExpression)
				if lhsIsSum && rhsIsLiteral {
					// Distribute literal to sum expression
					expr := NewAddExpression(
						NewMultiplyExpression(lhsSum.Lhs(), rhsLiteral),
						NewMultiplyExpression(lhsSum.Rhs(), rhsLiteral),
					)
					return expr.Simplify(inputs)
				}

				rhsSum, rhsIsSum := rhs.(*AddExpression)
				lhsLiteral, lhsIsLiteral := lhs.(*LiteralExpression)
				if rhsIsSum && lhsIsLiteral {
					// Distribute literal to sum expression
					expr := NewAddExpression(
						NewMultiplyExpression(rhsSum.Lhs(), lhsLiteral),
						NewMultiplyExpression(rhsSum.Rhs(), lhsLiteral),
					)
					return expr.Simplify(inputs)
				}

				if lhsMultiply, lhsIsMultiply := lhs.(*MultiplyExpression); lhsIsMultiply && rhsIsSingleValue {
					if lhsLiteral, lhsIsLiteral := lhsMultiply.Lhs().(*LiteralExpression); lhsIsLiteral {
						expr := NewMultiplyExpression(
							NewLiteralExpression(lhsLiteral.value*rhsSingleValue),
							lhsMultiply.Rhs(),
						)
						return expr.Simplify(inputs)
					}
					if rhsLiteral, rhsIsLiteral := lhsMultiply.Rhs().(*LiteralExpression); rhsIsLiteral {
						expr := NewMultiplyExpression(
							lhsMultiply.Lhs(),
							NewLiteralExpression(rhsLiteral.value*rhsSingleValue),
						)
						return expr.Simplify(inputs)
					}
				}

				if rhsMultiply, rhsIsMultiply := rhs.(*MultiplyExpression); rhsIsMultiply && lhsIsSingleValue {
					if lhsLiteral, lhsIsLiteral := rhsMultiply.Lhs().(*LiteralExpression); lhsIsLiteral {
						expr := NewMultiplyExpression(
							NewLiteralExpression(lhsLiteral.value*lhsSingleValue),
							rhsMultiply.Rhs(),
						)
						return expr.Simplify(inputs)
					}
					if rhsLiteral, rhsIsLiteral := rhsMultiply.Rhs().(*LiteralExpression); rhsIsLiteral {
						expr := NewMultiplyExpression(
							rhsMultiply.Lhs(),
							NewLiteralExpression(rhsLiteral.value*lhsSingleValue),
						)
						return expr.Simplify(inputs)
					}
				}
			*/
			return NewMultiplyExpression(lhs, rhs)
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

////////////////////////////////////////////////////////////////////////////////
// multiplyRange

func (r *multiplyRange) Includes(value int) bool {

	lhsBounded, lhsIsBounded := r.lhs.(BoundedRange)
	rhsBounded, rhsIsBounded := r.rhs.(BoundedRange)

	if lhsIsBounded && rhsIsBounded {
		inBounds := (value >= lhsBounded.Min()*rhsBounded.Min()) && (value <= lhsBounded.Max()*rhsBounded.Max())
		if !inBounds {
			return false
		}
	}

	nextValue := r.Values(fmt.Sprintf("%s includes %d", r, value))
	for v, ok := nextValue(); ok; v, ok = nextValue() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *multiplyRange) String() string {
	return fmt.Sprintf("<%s * %s>", r.lhs.String(), r.rhs.String())
}

func (r *multiplyRange) Values(context string) func() (int, bool) {

	pos := 0

	return func() (int, bool) {
		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue * rhsValue },
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
