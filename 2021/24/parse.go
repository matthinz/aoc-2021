package d24

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var expressionFactories = map[string]func(expressions ...interface{}) Expression{
	"add": NewAddExpression,
	"div": NewDivideExpression,
	"eql": NewEqualsExpression,
	"mod": NewModuloExpression,
	"mul": NewMultiplyExpression,
}

func parseInput(r io.Reader) *Registers {

	result := NewRegisters()

	inputIndex := 0

	s := bufio.NewScanner(r)
	lineIndex := 0

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		lineIndex++

		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")

		if len(parts) == 2 {

			// this must be an "inp"
			if parts[0] != "inp" {
				panic(fmt.Sprintf("Invalid unary operation: %s", parts[0]))
			}
			expr := NewInputExpression(inputIndex)

			inputIndex++

			result.set(parts[1], expr)

			continue
		}

		lhs := result.get(parts[1])

		var rhs Expression

		literalValue, err := strconv.ParseInt(parts[2], 10, 32)
		if err == nil {
			rhs = NewLiteralExpression(int(literalValue))
		} else {
			rhs = result.get(parts[2])
		}

		expr := makeBinaryExpression(parts[0], lhs, rhs)

		// set the value of the specified register to the expression
		// fmt.Printf("%d: %s \n", lineIndex, line)
		simplified := expr.Simplify([]int{})
		// fmt.Println(simplified.String())
		result.set(parts[1], simplified)

	}

	return result
}

func makeBinaryExpression(kind string, lhs Expression, rhs Expression) Expression {
	factory, ok := expressionFactories[kind]

	if !ok {
		panic(fmt.Sprintf("Invalid op: %s", kind))
	}

	return factory(lhs, rhs)
}
