package d24

import (
	"fmt"
	"log"
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
	lhs Range
	rhs Range
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

func (e *ModuloExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) % e.rhs.Evaluate(inputs)
}

func (e *ModuloExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	return findInputsForBinaryExpression(
		e,
		target,
		func(lhsValue int, rhsRange Range) (chan int, error) {
			// need to find rhsValues such that lhsValue % rhsValue = target
			// these would be factors of lhsValue - target
			// TODO: smarter way

			result := make(chan int)

			go func() {
				defer close(result)

				rhsValues := rhsRange.Values()
				for rhsValue := range rhsValues {
					if lhsValue%rhsValue == target {
						result <- rhsValue
					}
				}
			}()

			return result, nil
		},
		d,
		l,
	)
}

func (e *ModuloExpression) Range() Range {
	lhsRange := e.Lhs().Range()
	rhsRange := e.Rhs().Range()

	return &moduloRange{
		lhs: lhsRange,
		rhs: rhsRange,
	}
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

////////////////////////////////////////////////////////////////////////////////
// moduloRange

func (r *moduloRange) Includes(value int) bool {
	for v := range r.Values() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *moduloRange) Split(around int) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *moduloRange) Values() chan int {
	ch := make(chan int)

	go func() {

		lhsContinuous, lhsIsContinuous := r.lhs.(*continuousRange)
		rhsContinuous, rhsIsContinuous := r.rhs.(*continuousRange)

		var min, max *int
		var hitMin, hitMax int

		var prevValue *int

		for lhsValue := range r.lhs.Values() {
			for rhsValue := range r.rhs.Values() {

				value := lhsValue % rhsValue

				if min != nil && value == *min {
					hitMin++
				}

				if max != nil && value == *max {
					hitMax++
				}

				if hitMin > 1 && hitMax > 1 {
					continue
				}

				if prevValue != nil && value == *prevValue {
					continue
				}

				ch <- value
				prevValue = &value
			}
		}
		defer close(ch)
	}()

	return ch
}

func (r *moduloRange) String() string {
	const maxLength = 10
	values := make(map[int]bool)
	for value := range r.Values() {
		values[value] = true
		if len(values) > maxLength {
			return fmt.Sprintf("(%s %% %s)", r.lhs.String(), r.rhs.String())
		}
	}

	distinctValues := make([]int, 0)
	for value, _ := range values {
		distinctValues = append(distinctValues, value)
	}

	sort.Ints(distinctValues)

	stringValues := make([]string, 0)
	for value := range distinctValues {
		stringValues = append(stringValues, strconv.FormatInt(int64(value), 10))
	}

	return strings.Join(stringValues, ",")
}
