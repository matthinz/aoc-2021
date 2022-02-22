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
		if lhsContinuous.step == 1 {
			e.cachedRange = &continuousRange{
				min:  lhsContinuous.min + rhsContinuous.min,
				max:  lhsContinuous.max + rhsContinuous.max,
				step: rhsContinuous.step,
			}
		} else if rhsContinuous.step == 1 {
			e.cachedRange = &continuousRange{
				min:  lhsContinuous.min + rhsContinuous.min,
				max:  lhsContinuous.max + rhsContinuous.max,
				step: lhsContinuous.step,
			}
		}
	} else {
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

////////////////////////////////////////////////////////////////////////////////
// sumRange

func (r *sumRange) Includes(value int) bool {
	next := r.Values()
	for v, ok := next(); ok; v, ok = next() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *sumRange) Split(around Range) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *sumRange) String() string {
	return fmt.Sprintf("(%s + %s)", r.lhs.String(), r.rhs.String())
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
