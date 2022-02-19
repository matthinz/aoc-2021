package d24

import (
	"fmt"
	"log"
	"time"
)

// AddExpression defines a BinaryExpression that adds its left and righthand sides.
type AddExpression struct {
	binaryExpression
}

// sumRange is a Range implementation that represents two other Ranges
// summed together.
type sumRange struct {
	lhs, rhs Range
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

	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	if lhsIsContinuous && rhsIsContinuous {
		if lhsContinuous.step == 1 {
			return &continuousRange{
				min:  lhsContinuous.min + rhsContinuous.min,
				max:  lhsContinuous.max + rhsContinuous.max,
				step: rhsContinuous.step,
			}
		} else if rhsContinuous.step == 1 {
			return &continuousRange{
				min:  lhsContinuous.min + rhsContinuous.min,
				max:  lhsContinuous.max + rhsContinuous.max,
				step: lhsContinuous.step,
			}
		}
	}

	return &sumRange{
		lhs: lhsRange,
		rhs: rhsRange,
	}
}

func (e *AddExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	fmt.Printf("BEFORE lhs: %s\nrhs: %s\n", lhs.String(), rhs.String())

	start := time.Now()
	lhsRange := lhs.Range()
	fmt.Printf("AFTER lhs: %s (%s)\n", lhsRange.String(), time.Now().Sub(start))

	start = time.Now()
	rhsRange := rhs.Range()
	fmt.Printf("AFTER rhs: %s (%s)\n", rhsRange.String(), time.Now().Sub(start))

	// if both ranges are single numbers we are adding two literals
	lhsSingleValue, lhsIsSingleValue := GetSingleValueOfRange(lhsRange)
	rhsSingleValue, rhsIsSingleValue := GetSingleValueOfRange(rhsRange)

	if lhsIsSingleValue && rhsIsSingleValue {
		return NewLiteralExpression(lhsSingleValue + rhsSingleValue)
	}

	// if either range is zero, use the other
	if lhsIsSingleValue && lhsSingleValue == 0 {
		return rhs
	}

	if rhsIsSingleValue && rhsSingleValue == 0 {
		return lhs
	}

	result := NewAddExpression(lhs, rhs)
	result.(*AddExpression).isSimplified = true

	return result
}

////////////////////////////////////////////////////////////////////////////////
// sumRange

func (r *sumRange) Includes(value int) bool {
	values := r.Values()
	for v := range values {
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

func (r *sumRange) Values() chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		var prevValue *int
		var min *int
		var max *int

		for lhsValue := range r.lhs.Values() {
			for rhsValue := range r.rhs.Values() {

				value := lhsValue + rhsValue

				if prevValue != nil && value == *prevValue {
					continue
				}

				if min != nil && value == *min {
					continue
				}

				if max != nil && value == *max {
					continue
				}

				ch <- value

				prevValue = &value

				if min == nil || value < *min {
					min = &value
				}

				if max == nil || value > *max {
					max = &value
				}
			}
		}
	}()

	return ch
}
