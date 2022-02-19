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
					for rhsValue := range rhsRange.Values() {
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
	return &equalsRange{
		lhs: e.Lhs().Range(),
		rhs: e.Rhs().Range(),
	}
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

	values := r.Values()
	for v := range values {
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
	for value := range r.Values() {
		values = append(values, strconv.FormatInt(int64(value), 10))
	}

	return fmt.Sprintf("(%s)", strings.Join(values, ","))
}

func (r *equalsRange) Values() chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		// TODO: The actual values
		result <- 0
		result <- 1

	}()
	return result
}
