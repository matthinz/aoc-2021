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
	Simplify() Expression
	SimplifyUsingPartialInputs(inputs map[int]int) Expression
	String() string
}

type BinaryExpression interface {
	Lhs() Expression
	Rhs() Expression
}

// binaryExpression is an embeddable Expression comprised of two expressions,
// (left- and right-hand sides) and an operator.
type binaryExpression struct {
	lhs      Expression
	rhs      Expression
	operator rune

	cachedRange Range

	// Whether this expression has been simplified already.
	isSimplified bool
}

func (e *binaryExpression) Lhs() Expression {
	return e.lhs
}

func (e *binaryExpression) Rhs() Expression {
	return e.rhs
}

func (e *binaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.lhs.String(), string(e.operator), e.rhs.String())
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
