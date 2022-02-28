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
	Simplify(inputs map[int]int) Expression

	String() string
}

type BinaryExpression interface {
	Lhs() Expression
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

func evaluateBinaryExpression(e BinaryExpression, op func(lhs, rhs int) (int, error)) (int, error) {
	lhsValue, lhsError := e.Lhs().Evaluate()
	if lhsError != nil {
		return 0, lhsError
	}

	rhsValue, rhsError := e.Rhs().Evaluate()
	if rhsError != nil {
		return 0, rhsError
	}

	return op(lhsValue, rhsValue)
}

func simplifyBinaryExpression(
	e *binaryExpression,
	inputs map[int]int,
	simplifier func(lhs, rhs Expression) Expression,
) Expression {
	if len(inputs) == 0 && e.normalized != nil {
		// Simplifying with no inputs is the same as normalizing
		return e.normalized
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
	simplifiedRhs := ogRhs.Simplify(inputs)

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
