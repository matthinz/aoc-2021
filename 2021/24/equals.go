package d24

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type EqualsExpression struct {
	binaryExpression
}

type equalsRange struct {
	lhs Range
	rhs Range
}

func NewEqualsExpression(lhs, rhs Expression) Expression {
	return &EqualsExpression{
		binaryExpression: binaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '=',
		},
	}
}

func (e *EqualsExpression) Accept(visitor func(e Expression)) {
	visitor(e)
	e.lhs.Accept(visitor)
	e.rhs.Accept(visitor)
}

func (e *EqualsExpression) Evaluate(inputs []int) int {
	lhsValue := e.lhs.Evaluate(inputs)
	rhsValue := e.rhs.Evaluate(inputs)
	if lhsValue == rhsValue {
		return 1
	} else {
		return 0
	}
}

func (e *EqualsExpression) FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error) {
	if target != 0 && target != 1 {
		return nil, fmt.Errorf("EqualsExpression can't seek a target other than 0 or 1 (got %d)", target)
	}

	return findInputsForBinaryExpression(
		e,
		target,
		func(lhsValue int, rhsRange Range) (chan int, error) {
			ch := make(chan int)

			go func() {
				defer close(ch)

				if target == 0 {
					// We want any members of rhsRange *not equal to* lhsValue
					nextRhsValue := rhsRange.Values()

					for rhsValue, ok := nextRhsValue(); ok; rhsValue, ok = nextRhsValue() {
						if rhsValue != lhsValue {
							ch <- rhsValue
						}
					}

					return
				}

				// We must find a value in rhsRange that equals lhsValue
				if rhsRange.Includes(lhsValue) {
					ch <- lhsValue
				}
			}()

			return ch, nil
		},
		d,
		l,
	)
}

func (e *EqualsExpression) Range() Range {
	if e.cachedRange == nil {
		e.cachedRange = &equalsRange{
			lhs: e.Lhs().Range(),
			rhs: e.Rhs().Range(),
		}
	}
	return e.cachedRange
}

func (e *EqualsExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if all elements of both ranges are equal, we are comparing two equal values
	if RangesAreEqual(lhsRange, rhsRange) {
		return oneLiteral
	}

	// If the ranges of each side of the comparison will never intersect,
	// then we can always return "0" for this expression
	if !RangesIntersect(lhsRange, rhsRange) {
		return zeroLiteral
	}

	expr := NewEqualsExpression(lhs, rhs)
	expr.(*EqualsExpression).isSimplified = true

	return expr
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}

////////////////////////////////////////////////////////////////////////////////
// equalsRange

func (r *equalsRange) Includes(value int) bool {
	if value != 0 && value != 1 {
		return false
	}

	nextValue := r.Values()
	for v, ok := nextValue(); ok; v, ok = nextValue() {
		if v == value {
			return true
		}
	}

	return false
}

func (r *equalsRange) Split(around Range) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *equalsRange) String() string {
	var values []string
	nextValue := r.Values()
	for value, ok := nextValue(); ok; value, ok = nextValue() {
		values = append(values, strconv.FormatInt(int64(value), 10))
	}

	return fmt.Sprintf("(%s)", strings.Join(values, ","))
}

func (r *equalsRange) Values() func() (int, bool) {
	nextValue := 0
	return func() (int, bool) {

		if nextValue > 1 {
			return 0, false
		}

		value := nextValue
		nextValue++
		return value, true
	}
}
