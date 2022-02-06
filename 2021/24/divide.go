package d24

import (
	"fmt"
	"log"
)

type DivideExpression struct {
	binaryExpression
}

type divisionRange struct {
	lhs, rhs Range
}

func NewDivideExpression(lhs, rhs Expression) Expression {
	return &DivideExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '/',
		},
	}
}

func (e *DivideExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *DivideExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) / e.rhs.Evaluate(inputs)
}

func (e *DivideExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	return findInputsForBinaryExpression(
		e,
		target,
		func(dividend int, divisorRange Range) (chan int, error) {

			// 4th grade math recap: dividend / divisor = target
			// here we return potential divisors between min and max that will equal target

			ch := make(chan int)

			go func() {
				defer close(ch)

				if target == 0 {
					// When target == 0, divisor can't affect the result, except when it
					// can. We're doing integer division, so a large enough divisor *could* get us to zero
					// e.g if we're doing 6 / x = 0, any x > 6 will result in 0

					potentialDivisors := divisorRange.Values()

					if dividend == 0 {
						// It does not matter what the divisor is -- any non-zero value
						// in divisorRange will work
						for divisor := range potentialDivisors {
							if divisor != 0 {
								ch <- divisor
							}
						}
						return
					}

					// Any number *larger* than dividend will result in 0
					for divisor := range potentialDivisors {
						if divisor == 0 {
							continue
						}
						if divisor > dividend {
							ch <- divisor
						}
					}
					return
				}

				// dividend / divisor = target
				// dividend = divisor * target
				// divisor = dividend / target
				divisor := dividend / target

				if divisor == 0 {
					return
				}

				if dividend/divisor != target {
					return
				}

				if !divisorRange.Includes(divisor) {
					return
				}

				ch <- divisor
			}()

			return ch, nil
		},
		d,
		l,
	)
}

func (e *DivideExpression) Range() Range {
	return &divisionRange{
		lhs: e.Lhs().Range(),
		rhs: e.Rhs().Range(),
	}
}

func (e *DivideExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	lhsSingleValue, lhsIsSingleValue := GetSingleValueOfRange(lhsRange)
	rhsSingleValue, rhsIsSingleValue := GetSingleValueOfRange(rhsRange)

	if lhsIsSingleValue && rhsIsSingleValue {
		return NewLiteralExpression(lhsSingleValue / rhsSingleValue)
	}

	// if left value is zero, this will eval to zero
	if lhsIsSingleValue && lhsSingleValue == 0 {
		return NewLiteralExpression(0)
	}

	// if right value is 1, this will eval to lhs
	if rhsIsSingleValue && rhsSingleValue == 1 {
		return lhs
	}

	result := NewDivideExpression(lhs, rhs)
	result.(*DivideExpression).isSimplified = true

	return result
}

////////////////////////////////////////////////////////////////////////////////
// divisionRange

func (r *divisionRange) Includes(value int) bool {
	for i := range r.Values() {
		if i == value {
			return true
		}
	}
	return false
}

func (r *divisionRange) Split(around int) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *divisionRange) String() string {
	return fmt.Sprintf("%s + %s", r.lhs.String(), r.rhs.String())
}

func (r *divisionRange) Values() chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		var prevValue *int

		lhsValues := r.lhs.Values()
		for lhsValue := range lhsValues {
			rhsValues := r.rhs.Values()
			for rhsValue := range rhsValues {
				if rhsValue == 0 {
					continue
				}
				value := lhsValue / rhsValue
				if prevValue == nil || value != *prevValue {
					result <- value
					prevValue = &value
				}
			}
		}
	}()

	return result
}
