package d24

import (
	"fmt"
)

// Expression is a generic interface encapsulating an expression that can
// be evaluated with an ALU
type Expression interface {
	Accept(visitor func(e Expression))
	Evaluate(inputs []int) int
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
