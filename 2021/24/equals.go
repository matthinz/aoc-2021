package d24

import (
	"fmt"
	"log"
)

type EqualsExpression struct {
	binaryExpression
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
	return newContinuousRange(0, 1, 1)
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

func (e *EqualsExpression) SimplifyUsingPartialInputs(inputs map[int]int) Expression {
	lhs := e.Lhs().SimplifyUsingPartialInputs(inputs)
	rhs := e.Rhs().SimplifyUsingPartialInputs(inputs)
	expr := NewEqualsExpression(lhs, rhs)
	return expr.Simplify()
}

func (e *EqualsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}
