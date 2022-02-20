package d24

import (
	"fmt"
	"log"
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

func (e *MultiplyExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	return findInputsForBinaryExpression(
		e,
		target,
		func(lhsValue int, rhsRange Range) (chan int, error) {

			ch := make(chan int)

			go func() {
				defer close(ch)

				if target == 0 {

					if lhsValue != 0 {
						// rhsValue *must* be zero
						if rhsRange.Includes(0) {
							ch <- 0
							return
						}
					}
				} else if target == lhsValue {
					if rhsRange.Includes(1) {
						ch <- 1
					}
					return
				}

				for rhsValue := range rhsRange.Values() {
					if lhsValue*rhsValue == target {
						ch <- rhsValue
					}
				}
			}()

			return ch, nil
		},
		d,
		l,
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
			e.cachedRange = newContinuousRange(
				lhsContinuous.min*rhsContinuous.min,
				lhsContinuous.max*rhsContinuous.max,
				lhsContinuous.min,
			)
		} else if rhsContinuous.min == rhsContinuous.max {
			e.cachedRange = newContinuousRange(
				lhsContinuous.min*rhsContinuous.min,
				lhsContinuous.max*rhsContinuous.max,
				rhsContinuous.min,
			)
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

	expr := NewMultiplyExpression(lhs, rhs)
	expr.(*MultiplyExpression).isSimplified = true

	return expr
}

////////////////////////////////////////////////////////////////////////////////
// multiplyRange

func (r *multiplyRange) Includes(value int) bool {
	for i := range r.Values() {
		if i == value {
			return true
		}
	}
	return false
}

func (r *multiplyRange) Split(around Range) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *multiplyRange) String() string {
	return fmt.Sprintf("(%s * %s)", r.lhs.String(), r.rhs.String())
}

func (r *multiplyRange) Values() chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue * rhsValue },
			)
		}

		for _, value := range *r.cachedValues {
			ch <- value
		}
	}()

	return ch
}
