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

		lhsRange := e.lhs.Range()
		rhsRange := e.rhs.Range()

		lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
		rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

		if lhsIsContinuous && rhsIsContinuous {

			if rhsContinuous.min == rhsContinuous.max {

				rhsIsFactorOfLhsStep := (lhsContinuous.step/rhsContinuous.min)*rhsContinuous.min == lhsContinuous.step

				if rhsIsFactorOfLhsStep {
					// If lhs is continuous and rhs is a factor of the step of lhs,
					// then we can cleanly divide everything
					e.cachedRange = newContinuousRange(
						lhsContinuous.min/rhsContinuous.min,
						lhsContinuous.max/rhsContinuous.max,
						lhsContinuous.step/rhsContinuous.min,
					)
				}

			}

		}

		if e.cachedRange == nil {

			e.cachedRange = &divisionRange{
				lhs: e.Lhs().Range(),
				rhs: e.Rhs().Range(),
			}
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

	lhsContinuous, lhsIsContinuous := lhsRange.(*continuousRange)
	rhsContinuous, rhsIsContinuous := rhsRange.(*continuousRange)

	if lhsIsContinuous && lhsContinuous.min == 0 && lhsContinuous.max == 0 {
		return NewLiteralExpression(0)
	}

	if lhsIsContinuous && rhsIsContinuous {
		if lhsContinuous.max < rhsContinuous.min {
			// this will _always_ be zero
			return NewLiteralExpression(0)
		}
	}

	if rhsIsContinuous && rhsContinuous.min == rhsContinuous.max {
		expr := divideExpressionByInt(lhs, rhsContinuous.min)
		if expr != nil {
			return expr.Simplify()
		}
	}

	result := NewDivideExpression(lhs, rhs)
	result.(*DivideExpression).isSimplified = true

	return result
}

// Attempts to divide an expression by an integer value, returning
// a new Expression if successful. If the operation is not possible, returns
// nil.
func divideExpressionByInt(dividend Expression, divisor int) Expression {
	if divisor == 0 {
		return nil
	}

	if divisor == 1 {
		return dividend
	}

	r := dividend.Range()

	bounds := getBounds(r)
	if bounds != nil {
		// If the range of our dividend is completely less than the divisor,
		// we can just zero the whole thing out
		if bounds.Max() < divisor {
			return NewLiteralExpression(0)
		}
	}

	if literal, isLiteral := dividend.(*LiteralExpression); isLiteral {
		return divideLiteralByInt(*literal, divisor)
	}

	input, isInput := dividend.(*InputExpression)
	if isInput {
		expr := NewDivideExpression(input, NewLiteralExpression(divisor))
		expr.(*DivideExpression).isSimplified = true
		return expr
	}

	if sum, isSum := dividend.(*AddExpression); isSum {
		lhs := divideExpressionByInt(sum.lhs, divisor)
		rhs := divideExpressionByInt(sum.rhs, divisor)
		if lhs != nil && rhs != nil {
			fmt.Printf("lhs: %s, rhs: %s\n", lhs.String(), rhs.String())
			return NewAddExpression(lhs, rhs)
		}
	}

	if multiply, isMultiply := dividend.(*MultiplyExpression); isMultiply {
		return divideMultiplyByInt(*multiply, divisor)
	}

	return nil
}

func divideLiteralByInt(dividend LiteralExpression, divisor int) Expression {
	// We avoid doing things that could lose precision
	isSafe := dividend.value%divisor == 0
	if isSafe {
		return NewLiteralExpression(dividend.value / divisor)
	} else {
		return nil
	}
}

func divideMultiplyByInt(dividend MultiplyExpression, divisor int) Expression {

	lhsLiteral, lhsIsLiteral := dividend.lhs.(*LiteralExpression)
	rhsLiteral, rhsIsLiteral := dividend.rhs.(*LiteralExpression)

	if lhsIsLiteral {
		lhs := divideExpressionByInt(lhsLiteral, divisor)
		if lhs != nil {
			return NewMultiplyExpression(lhs, dividend.rhs)
		}
	}

	if rhsIsLiteral {
		rhs := divideExpressionByInt(rhsLiteral, divisor)
		if rhs != nil {
			return NewMultiplyExpression(dividend.lhs, rhs)
		}
	}

	// If either side of the dividend is itself a MultiplyExpression, attempt
	// to find a subexpression we can cleanly apply the divisor to
	if lhsMultiply, lhsIsMultiply := dividend.lhs.(*MultiplyExpression); lhsIsMultiply {
		if newLhs := divideMultiplyByInt(*lhsMultiply, divisor); newLhs != nil {
			return NewMultiplyExpression(
				newLhs,
				dividend.Rhs(),
			)
		}
	}

	if rhsMultiply, rhsIsMultiply := dividend.rhs.(*MultiplyExpression); rhsIsMultiply {
		if newRhs := divideMultiplyByInt(*rhsMultiply, divisor); newRhs != nil {
			return NewMultiplyExpression(
				dividend.Lhs(),
				newRhs,
			)
		}
	}

	return nil
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

func (r *divisionRange) String() string {
	return fmt.Sprintf("<%s / %s>", r.lhs.String(), r.rhs.String())
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
