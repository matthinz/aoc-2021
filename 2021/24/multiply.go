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

func NewMultiplyExpression(lhs, rhs Expression) Expression {
	return &MultiplyExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '*',
		},
	}
}

func (e *MultiplyExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *MultiplyExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) * e.rhs.Evaluate(inputs)
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

func (e *MultiplyExpression) Simplify() Expression {
	if e.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

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
		return expr.Simplify()
	}

	rhsSum, rhsIsSum := rhs.(*AddExpression)
	lhsLiteral, lhsIsLiteral := lhs.(*LiteralExpression)
	if rhsIsSum && lhsIsLiteral {
		// Distribute literal to sum expression
		expr := NewAddExpression(
			NewMultiplyExpression(rhsSum.Lhs(), lhsLiteral),
			NewMultiplyExpression(rhsSum.Rhs(), lhsLiteral),
		)
		return expr.Simplify()
	}

	if lhsMultiply, lhsIsMultiply := lhs.(*MultiplyExpression); lhsIsMultiply && rhsIsSingleValue {
		if lhsLiteral, lhsIsLiteral := lhsMultiply.Lhs().(*LiteralExpression); lhsIsLiteral {
			expr := NewMultiplyExpression(
				NewLiteralExpression(lhsLiteral.value*rhsSingleValue),
				lhsMultiply.Rhs(),
			)
			return expr.Simplify()
		}
		if rhsLiteral, rhsIsLiteral := lhsMultiply.Rhs().(*LiteralExpression); rhsIsLiteral {
			expr := NewMultiplyExpression(
				lhsMultiply.Lhs(),
				NewLiteralExpression(rhsLiteral.value*rhsSingleValue),
			)
			return expr.Simplify()
		}
	}

	if rhsMultiply, rhsIsMultiply := rhs.(*MultiplyExpression); rhsIsMultiply && lhsIsSingleValue {
		if lhsLiteral, lhsIsLiteral := rhsMultiply.Lhs().(*LiteralExpression); lhsIsLiteral {
			expr := NewMultiplyExpression(
				NewLiteralExpression(lhsLiteral.value*lhsSingleValue),
				rhsMultiply.Rhs(),
			)
			return expr.Simplify()
		}
		if rhsLiteral, rhsIsLiteral := rhsMultiply.Rhs().(*LiteralExpression); rhsIsLiteral {
			expr := NewMultiplyExpression(
				rhsMultiply.Lhs(),
				NewLiteralExpression(rhsLiteral.value*lhsSingleValue),
			)
			return expr.Simplify()
		}
	}

	expr := NewMultiplyExpression(lhs, rhs)
	expr.(*MultiplyExpression).isSimplified = true

	return expr
}

func (e *MultiplyExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewMultiplyExpression(lhs, rhs)
	return expr.Simplify()
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

	nextValue := r.Values()
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

func (r *multiplyRange) Values() func() (int, bool) {

	pos := 0

	return func() (int, bool) {
		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue * rhsValue },
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
