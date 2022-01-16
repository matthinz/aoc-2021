package d24

import (
	"math"
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
	BinaryExpression
}

func NewModuloExpression(lhs, rhs Expression) Expression {
	return &ModuloExpression{
		BinaryExpression: BinaryExpression{
			lhs:      lhs,
			rhs:      rhs,
			operator: '%',
		},
	}
}

func (e *ModuloExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) % e.rhs.Evaluate(inputs)
}

func (e *ModuloExpression) FindInputs(target int, d InputDecider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.BinaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			// need to find rhsValues such that lhsValue % rhsValue = target
			// these would be factors of lhsValue - target
			// TODO: smarter way

			var result []int
			for i := rhsRange.min; i <= rhsRange.max; i++ {
				if lhsValue%i == target {
					result = append(result, i)
				}
			}

			return result, nil
		},
		d,
	)
}

func (e *ModuloExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	// 1. Probe for the upper and lower limits of the range
	//    This allows us to know, when iterating over large sets of numbers,
	//    when we should stop iterating.

	lowerBoundary, upperBoundary := findModuloBoundaries(lhsRange, rhsRange)

	// 2. Iterate over all possible combinations of values and find the
	//    minimum and maximum results. Once our min and max values are
	//    equal to our boundaries, we know there's no more we can do.

	min := math.MaxInt
	max := math.MinInt

	lhsRange.Each(func(i int) bool {
		rhsRange.Each(func(j int) bool {

			value := i % j

			if value < min {
				min = value
			}

			if value > max {
				max = value
			}

			return min != lowerBoundary || max != upperBoundary
		})

		return min != lowerBoundary || max != upperBoundary
	})

	return NewIntRange(min, max)
}

func (e *ModuloExpression) Simplify() Expression {
	if e.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if lhs is 0, we can resolve to zero
	if lhsRange.EqualsInt(0) {
		return zeroLiteral
	}

	// if both ranges are single numbers, we can evaluate to a literal
	if lhsRange.Len() == 1 && rhsRange.Len() == 1 {
		return NewLiteralExpression(lhsRange.min % rhsRange.min)
	}

	// if lhs is 1 number and *less than* the rhs range, we can eval to a literal
	if lhsRange.Len() == 1 && lhsRange.min < rhsRange.min {
		return NewLiteralExpression(lhsRange.min)
	}

	expr := NewModuloExpression(lhs, rhs)
	expr.(*ModuloExpression).isSimplified = true
	return expr
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Finds the maximum and minimum *possible* modulo values
func findModuloBoundaries(lhsRange, rhsRange IntRange) (int, int) {

	maxInt := func(values ...int) int {
		result := math.MinInt
		for _, value := range values {
			if value > result {
				result = value
			}
		}
		return result
	}

	minInt := func(values ...int) int {
		result := math.MaxInt
		for _, value := range values {
			if value < result {
				result = value
			}
		}
		return result
	}

	lowerBound := math.MaxInt
	upperBound := math.MinInt

	lhsNegative, lhsZero, lhsPositive := lhsRange.Split(0)
	rhsNegative, _, rhsPositive := rhsRange.Split(0)

	if lhsNegative != nil {
		// Having negative values on the left hand side of the modulo operation
		// means that it is possible the result could be negative.
		lowerBound = minInt(lowerBound, 0)

		if rhsNegative != nil {
			lowerBound = minInt(lowerBound, rhsNegative.min+1)
		}

		if rhsPositive != nil {
			lowerBound = minInt(lowerBound, rhsPositive.max*-1)
		}

		upperBound = maxInt(upperBound, lowerBound)
	}

	if lhsZero != nil {
		// When LHS is zero, the result will be zero
		lowerBound = minInt(lowerBound, 0)
		upperBound = maxInt(lowerBound, 0)
	}

	if lhsPositive != nil {
		// When LHS is positive, the result will be positive
		lowerBound = minInt(lowerBound, 0)

		if rhsNegative != nil {
			lowerBound = minInt(lowerBound, rhsNegative.max*-1)
			upperBound = maxInt(upperBound, rhsNegative.min*-1)
		}

		if rhsPositive != nil {
			lowerBound = minInt(lowerBound, rhsPositive.min-1)
			upperBound = maxInt(upperBound, rhsPositive.max-1)
		}
	}

	return lowerBound, upperBound
}

func (e *ModuloExpression) Visit(v func(e Expression)) {
	v(e)
	e.lhs.Visit(v)
	e.rhs.Visit(v)
}
