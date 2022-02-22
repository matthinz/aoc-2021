package d24

import (
	"fmt"
	"log"
)

type DivideExpression struct {
	binaryExpression
}

type divisionRange struct {
	lhs, rhs     Range
	cachedValues *[]int
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

					nextPotentialDivisor := divisorRange.Values()

					if dividend == 0 {
						// It does not matter what the divisor is -- any non-zero value
						// in divisorRange will work
						for divisor, ok := nextPotentialDivisor(); ok; divisor, ok = nextPotentialDivisor() {
							if divisor != 0 {
								ch <- divisor
							}
						}
						return
					}

					// Any number *larger* than dividend will result in 0
					for divisor, ok := nextPotentialDivisor(); ok; divisor, ok = nextPotentialDivisor() {
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
	if e.cachedRange == nil {
		e.cachedRange = &divisionRange{
			lhs: e.Lhs().Range(),
			rhs: e.Rhs().Range(),
		}
	}

	return e.cachedRange
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
	next := r.Values()

	for v, ok := next(); ok; v, ok = next() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *divisionRange) Split(around Range) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *divisionRange) String() string {
	return fmt.Sprintf("%s + %s", r.lhs.String(), r.rhs.String())
}

func (r *divisionRange) Values() func() (int, bool) {

	pos := 0

	return func() (int, bool) {

		if r.cachedValues == nil {
			r.cachedValues = buildBinaryExpressionRangeValues(
				r.lhs,
				r.rhs,
				func(lhsValue, rhsValue int) int { return lhsValue / rhsValue },
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
