package d24

import (
	"fmt"
)

// Expression is a generic interface encapsulating an expression that can
// be evaluated with an ALU
type Expression interface {
	Accept(visitor func(e Expression))

	// Evaluates this expression and returns an integer value or an error if
	// evaluation fails.
	Evaluate() (int, error)

	// Returns the range of output values for this expression
	Range() Range

	// Given a set of known inputs, attempts to simplify this Expression.
	// Returns the simplified version.
	Simplify(inputs []int) Expression

	String() string
}

type BinaryExpression interface {
	Expression
	Lhs() Expression
	Operator() string
	Rhs() Expression
}

type Normalizer interface {
	MarkNormalized(expr Expression)
}

// binaryExpression is an embeddable Expression comprised of two expressions,
// (left- and right-hand sides) and an operator.
type binaryExpression struct {
	lhs      Expression
	rhs      Expression
	operator rune

	cachedRange Range

	// The indices of any inputs referenced by this expression or its sub-expressions
	referencedInputs *[]int

	cachedEvaluation *struct {
		value int
		err   error
	}

	// A cached normalized version of this expression
	normalized Expression
}

func (e *binaryExpression) Lhs() Expression {
	return e.lhs
}

func (e *binaryExpression) MarkNormalized(expr Expression) {
	e.normalized = expr
}

func (e *binaryExpression) Rhs() Expression {
	return e.rhs
}

func (e *binaryExpression) Operator() string {
	return string(e.operator)
}

func (e *binaryExpression) ReferencedInputs() []int {
	if e.referencedInputs != nil {
		return *e.referencedInputs
	}

	inputs := make(map[int]bool)
	e.Lhs().Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputs[input.index] = true
		}
	})
	e.Rhs().Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputs[input.index] = true
		}
	})

	result := make([]int, 0, len(inputs))
	for index := range inputs {
		result = append(result, index)
	}

	e.referencedInputs = &result

	return result
}

func (e *binaryExpression) String() string {
	var lhs, rhs string

	if e.lhs == nil {
		lhs = "<nil>"
	} else {
		lhs = e.lhs.String()
	}

	if e.rhs == nil {
		rhs = "<nil>"
	} else {
		rhs = e.rhs.String()
	}

	return fmt.Sprintf("(%s %s %s)", lhs, string(e.operator), rhs)
}

func newBinaryExpression(
	operator rune,
	factory func(expressions ...interface{}) Expression,
	expressions []interface{},
) binaryExpression {
	var result = binaryExpression{
		operator: operator,
		lhs:      nil,
		rhs:      nil,
	}

	for i := range expressions {
		if expressions[i] == nil {
			continue
		}

		var expr Expression

		// XXX: This is presumably faster than reflection?
		switch value := expressions[i].(type) {
		case []*AddExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []*DivideExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []*InputExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []*LiteralExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []*ModuloExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []*MultiplyExpression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case []Expression:
			for j := range value {
				if value[j] == nil {
					continue
				}
				if expr == nil {
					expr = value[j]
				} else if value[j] != nil {
					expr = factory(expr, value[j])
				}
			}
		case int:
			expr = NewLiteralExpression(value)
		default:
			expr = expressions[i].(Expression)
		}

		if expr == nil {
			continue
		}

		if result.lhs == nil {
			result.lhs = expr
		} else if result.rhs == nil {
			result.rhs = expr
		} else {
			newLhs := factory(result.lhs, result.rhs)
			result.lhs = newLhs
			result.rhs = expr
		}
	}

	return result
}

func buildBinaryExpressionRange(
	context string,
	e *binaryExpression,
	computeSingleValue func(lhs, rhs int) (int, error),
	computeSingleLhsValue func(lhs int, rhs ContinuousRange) (Range, error),
	computeSingleRhsValue func(lhs ContinuousRange, rhs int) (Range, error),
	compute func(lhs, rhs ContinuousRange) (Range, error),
) Range {

	if e.cachedRange != nil {
		return e.cachedRange
	}

	if compute == nil {
		compute = func(lhs, rhs ContinuousRange) (Range, error) {
			var ranges []ContinuousRange

			if lhs.Length() < rhs.Length() {
				nextLhs := lhs.Values("buildBinaryExpressionRange")

				for lhs, lhsOk := nextLhs(); lhsOk; lhs, lhsOk = nextLhs() {
					r, err := computeSingleLhsValue(lhs, rhs)
					if err == nil {
						ranges = append(ranges, r.Split()...)
					}
				}
			} else {
				nextRhs := rhs.Values("buildBinaryExpressionRange")
				for rhs, rhsOk := nextRhs(); rhsOk; rhs, rhsOk = nextRhs() {
					r, err := computeSingleRhsValue(lhs, rhs)
					if err == nil {
						ranges = append(ranges, r.Split()...)
					}
				}
			}

			if len(ranges) == 0 {
				return EmptyRange, nil
			}

			result := NewContinuousRangeSet(ranges)
			return result, nil
		}
	}

	if computeSingleRhsValue == nil {
		computeSingleRhsValue = func(lhs ContinuousRange, rhs int) (Range, error) {
			return computeSingleLhsValue(rhs, lhs)
		}
	}

	lhs, rhs := e.Lhs(), e.Rhs()
	lhsRanges, rhsRanges := lhs.Range().Split(), rhs.Range().Split()

	if len(lhsRanges) == 0 || len(rhsRanges) == 0 {
		return EmptyRange
	}

	var rawRanges []ContinuousRange
	singleValues := make(map[int]bool)

	for _, l := range lhsRanges {
		for _, r := range rhsRanges {
			if l.Min() == l.Max() {
				if r.Min() == r.Max() {
					value, err := computeSingleValue(l.Min(), r.Min())
					if err == nil {
						singleValues[value] = true
					}
				} else {
					value, err := computeSingleLhsValue(l.Min(), r)
					if err == nil {
						rawRanges = append(rawRanges, value.Split()...)
					}
				}
			} else {
				if r.Min() == r.Max() {
					value, err := computeSingleRhsValue(l, r.Min())
					if err == nil {
						rawRanges = append(rawRanges, value.Split()...)
					}
				} else {
					value, err := compute(l, r)
					if err == nil {
						rawRanges = append(rawRanges, value.Split()...)
					}
				}
			}
		}
	}

	resultRanges := make([]ContinuousRange, 0)

	for _, r := range rawRanges {
		if r.Length() < 10000000 {
			next := r.Values("buildBinaryExpressionRange")
			for value, ok := next(); ok; value, ok = next() {
				singleValues[value] = true
			}
		} else {
			resultRanges = append(resultRanges, r)
		}
	}

	if len(singleValues) > 0 {
		resultRanges = append(resultRanges, newRangeFromInts(singleValues).Split()...)
	}

	e.cachedRange = NewContinuousRangeSet(resultRanges)
	return e.cachedRange
}

func evaluateBinaryExpression(e BinaryExpression, op func(lhs, rhs int) (int, error)) (int, error) {
	switch x := e.(type) {
	case *AddExpression:
		if x.binaryExpression.cachedEvaluation != nil {
			return x.binaryExpression.cachedEvaluation.value, x.binaryExpression.cachedEvaluation.err
		}
	case *DivideExpression:
		if x.binaryExpression.cachedEvaluation != nil {
			return x.binaryExpression.cachedEvaluation.value, x.binaryExpression.cachedEvaluation.err
		}
	case *EqualsExpression:
		if x.binaryExpression.cachedEvaluation != nil {
			return x.binaryExpression.cachedEvaluation.value, x.binaryExpression.cachedEvaluation.err
		}
	case *ModuloExpression:
		if x.binaryExpression.cachedEvaluation != nil {
			return x.binaryExpression.cachedEvaluation.value, x.binaryExpression.cachedEvaluation.err
		}
	case *MultiplyExpression:
		if x.binaryExpression.cachedEvaluation != nil {
			return x.binaryExpression.cachedEvaluation.value, x.binaryExpression.cachedEvaluation.err
		}

	}

	lhsValue, lhsError := e.Lhs().Evaluate()
	if lhsError != nil {
		return 0, lhsError
	}

	rhsValue, rhsError := e.Rhs().Evaluate()
	if rhsError != nil {
		return 0, rhsError
	}

	result, err := op(lhsValue, rhsValue)

	switch x := e.(type) {
	case *AddExpression:
		x.binaryExpression.cachedEvaluation = &struct {
			value int
			err   error
		}{result, err}
	case *DivideExpression:
		x.binaryExpression.cachedEvaluation = &struct {
			value int
			err   error
		}{result, err}
	case *EqualsExpression:
		x.binaryExpression.cachedEvaluation = &struct {
			value int
			err   error
		}{result, err}
	case *ModuloExpression:
		x.binaryExpression.cachedEvaluation = &struct {
			value int
			err   error
		}{result, err}
	case *MultiplyExpression:
		x.binaryExpression.cachedEvaluation = &struct {
			value int
			err   error
		}{result, err}
	}

	return result, err

}

func simplifyBinaryExpression(
	e *binaryExpression,
	inputs []int,
	simplifier func(lhs, rhs Expression) Expression,
) Expression {
	if e.normalized != nil {
		if len(inputs) == 0 {
			// Simplifying with no inputs is the same as normalizing
			return e.normalized
		}
	}

	ogLhs := e.Lhs()
	if ogLhs == nil {
		panic(fmt.Sprintf("Lhs() returned nil on %v", e))
	}

	ogRhs := e.Rhs()
	if ogRhs == nil {
		panic(fmt.Sprintf("Rhs() returned nil on %v", e))
	}

	simplifiedLhs := ogLhs.Simplify(inputs)
	if simplifiedLhs == nil {
		simplifiedLhs = ogLhs
	}

	simplifiedRhs := ogRhs.Simplify(inputs)
	if simplifiedRhs == nil {
		simplifiedRhs = ogRhs
	}

	simplified := simplifier(simplifiedLhs, simplifiedRhs)

	if len(inputs) == 0 {
		e.normalized = simplified

		// Tell our simplified expression it is the normalized version of itself
		if n, isNormalizer := simplified.(Normalizer); isNormalizer {
			n.MarkNormalized(simplified)
		}
	}

	return simplified
}
