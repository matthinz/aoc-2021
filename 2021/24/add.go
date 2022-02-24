package d24

import (
	"fmt"
	"log"
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

func (e *AddExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) + e.rhs.Evaluate(inputs)
}

func (e *AddExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	return findInputsForBinaryExpression(
		e,
		target,
		func(lhsValue int, rhsRange Range) (chan int, error) {
			ch := make(chan int)

			go func() {
				defer close(ch)

				rhsValue := target - lhsValue
				if rhsRange.Includes(rhsValue) {
					ch <- rhsValue
				}
			}()

			return ch, nil
		},
		d,
		l,
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

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	lhsIsSingleValue := lhsIsContinuous && lhsContinuous.min == lhsContinuous.max
	rhsIsSingleValue := rhsIsContinuous && rhsContinuous.min == rhsContinuous.max

	if lhsIsSingleValue && rhsIsSingleValue {
		return NewLiteralExpression(lhsContinuous.min + rhsContinuous.min)
	} else if lhsIsSingleValue && lhsContinuous.min == 0 {
		return rhs
	} else if rhsIsSingleValue && rhsContinuous.min == 0 {
		return lhs
	}

	result := NewAddExpression(lhs, rhs)
	result.(*AddExpression).isSimplified = true

	return result
}

func (e *AddExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewAddExpression(lhs, rhs)
	return expr.Simplify()
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

	next := r.Values()
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

func (r *sumRange) Values() func() (int, bool) {

	pos := 0

	return func() (int, bool) {
		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue + rhsValue },
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
