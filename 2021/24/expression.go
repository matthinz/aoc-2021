package d24

import (
	"fmt"
	"log"
	"math"
)

// InputDecider takes two sets of inputs and decides which one to use.
// If no decision can be made, it should return an error.
type InputDecider func(a, b map[int]int) (map[int]int, error)

// Expression is a generic interface encapsulating an expression that can
// be evaluated with an ALU
type Expression interface {
	Accept(visitor func(e Expression))
	Evaluate(inputs []int) int
	// Returns a set of inputs that will make this expression evaluate to <target>.
	// <d> is a function that, given two potential sets of inputs, returns the one that should be preferred.
	FindInputs(target int, d InputDecider, l *log.Logger) (map[int]int, error)
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

func PreferFirstSetOfInputs(a, b map[int]int) (map[int]int, error) {
	return a, nil
}

func PreferSecondSetOfInputs(a, b map[int]int) (map[int]int, error) {
	return a, nil
}

func PreferInputsThatMakeLargerNumber(a, b map[int]int) (map[int]int, error) {

	aValue := inputsToNumber(a)
	bValue := inputsToNumber(b)

	if aValue >= bValue {
		return a, nil
	} else {
		return b, nil
	}
}

func inputsToNumber(inputs map[int]int) int {
	var digits []int

	for inputIndex, inputValue := range inputs {
		if len(digits) < inputIndex+1 {
			temp := make([]int, inputIndex+1)
			copy(temp, digits)
			digits = temp
		}
		digits[inputIndex] = inputValue
	}

	var result int

	for i, digit := range digits {
		power := (len(digits) - (i + 1))
		result += int(math.Pow10(power)) * digit
	}

	return result
}
