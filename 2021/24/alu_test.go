package d24

import (
	"fmt"
	"strings"
	"testing"
)

func TestFindAllInputsInZ(t *testing.T) {

	r := parseInput(strings.NewReader(realInput))

	r.z.Accept(func(expr Expression) {

		if be, ok := expr.(BinaryExpression); ok {

			lhs := be.Lhs()
			rhs := be.Rhs()

			_, lhsIsInput := lhs.(*InputExpression)
			_, rhsIsInput := rhs.(*InputExpression)

			if lhsIsInput || rhsIsInput {
				fmt.Printf("%s --- range: %v\n", be.String(), expr.Range())
			}

		}

	})

}
