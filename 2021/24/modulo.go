package d24

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Things about modulus
//
// x % y = "remainder of x/y"
//
// This works out to: (decimal part of x / y) * y
//
// When x < y, there is *only* a decimal part, so this becomes (x/y*y) == x
//
// When x == y, then modulus == 0
//
// When x > y:
//
// 	0 <= result <= y - 1
//
// When x == 0, result is *always* 0
//
// When x < 0 and y >= 0 result is negative
// When x >= 0 and y < 0 result is positive
// When x < 0 and y < 0 result is negative

type ModuloExpression struct {
	binaryExpression
}

type moduloRange struct {
	lhs          Range
	rhs          Range
	cachedValues *[]int
}

func NewModuloExpression(lhs, rhs Expression) Expression {
	return &ModuloExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '%',
		},
	}
}

func (e *ModuloExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *ModuloExpression) Evaluate() (int, error) {
	return evaluateBinaryExpression(
		e,
		func(lhs, rhs int) (int, error) {
			if rhs == 0 {
				return 0, fmt.Errorf("Cannot take modulo 0")
			}
			return lhs % rhs, nil
		},
	)
}

func (e *ModuloExpression) Range() Range {
	if e.cachedRange == nil {

		lhsRange := e.Lhs().Range()
		rhsRange := e.Rhs().Range()

		lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
		rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

		if lhsIsContinuous && rhsIsContinuous {
			if rhsContinuous.min == rhsContinuous.max {
				if lhsContinuous.step == rhsContinuous.min {
					// When the step of the lhs range == value of rhs range, then
					// then the lhs range is reduced to a single value.
					value := lhsContinuous.min % rhsContinuous.min
					e.cachedRange = newContinuousRange(value, value, 1)
				}
			}
		}

		if e.cachedRange == nil {
			e.cachedRange = &moduloRange{
				lhs: lhsRange,
				rhs: rhsRange,
			}
		}
	}

	return e.cachedRange
}

func (e *ModuloExpression) Simplify() Expression {
	if e.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	lhsSingleValue, lhsIsSingleValue := GetSingleValueOfRange(lhsRange)
	rhsSingleValue, rhsIsSingleValue := GetSingleValueOfRange(rhsRange)

	// If lhs is 0, we can resolve to zero
	if lhsIsSingleValue && lhsSingleValue == 0 {
		return zeroLiteral
	}

	// If both ranges are single numbers, we can simplify to a literal
	if lhsIsSingleValue && rhsIsSingleValue {
		return NewLiteralExpression(lhsSingleValue % rhsSingleValue)
	}

	// TODO: If lhs is 1 number and *less than* the rhs range, we can eval to a literal

	expr := NewModuloExpression(lhs, rhs)
	expr.(*ModuloExpression).isSimplified = true
	return expr
}

func (e *ModuloExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewModuloExpression(lhs, rhs)
	return expr.Simplify()
}

////////////////////////////////////////////////////////////////////////////////
// moduloRange

func (r *moduloRange) Includes(value int) bool {
	nextValue := r.Values(fmt.Sprintf("%s includes %d", r, value))
	for v, ok := nextValue(); ok; v, ok = nextValue() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *moduloRange) Values(context string) func() (int, bool) {

	pos := 0

	return func() (int, bool) {

		if r.cachedValues == nil {

			uniqueValues := make([]int, 0)

			nextRhsValue := r.rhs.Values(context)
			for rhsValue, ok := nextRhsValue(); ok; rhsValue, ok = nextRhsValue() {
				values := make(map[int]int)
				nextLhsValue := r.lhs.Values(context)

				for lhsValue, ok := nextLhsValue(); ok; lhsValue, ok = nextLhsValue() {
					value := lhsValue % rhsValue
					values[value]++

					if values[0] >= 2 {
						break
					}
				}

				for value := range values {
					uniqueValues = append(uniqueValues, value)
				}
			}
			r.cachedValues = &uniqueValues
		}

		if pos >= len(*r.cachedValues) {
			return 0, false
		}

		value := (*r.cachedValues)[pos]
		pos++
		return value, true
	}
}

func (r *moduloRange) String() string {
	const maxLength = 10
	values := make(map[int]bool)

	nextValue := r.Values("moduloRange.String()")
	for value, ok := nextValue(); ok; value, ok = nextValue() {
		values[value] = true
		if len(values) > maxLength {
			return fmt.Sprintf("<%s %% %s>", r.lhs.String(), r.rhs.String())
		}
	}

	distinctValues := make([]int, 0)
	for value := range values {
		distinctValues = append(distinctValues, value)
	}

	sort.Ints(distinctValues)

	stringValues := make([]string, 0)
	for value := range distinctValues {
		stringValues = append(stringValues, strconv.FormatInt(int64(value), 10))
	}

	return strings.Join(stringValues, ",")
}

func makeContinuous(a, b, c Range) (*continuousRange, *continuousRange, *continuousRange) {
	var outA, outB, outC *continuousRange

	if a != nil {
		outA = a.(*continuousRange)
	}

	if b != nil {
		outB = b.(*continuousRange)
	}

	if c != nil {
		outC = c.(*continuousRange)
	}

	return outA, outB, outC

}
