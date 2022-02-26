package d24

import (
	"sort"
	"testing"
)

func TestDivideExpressionEvaluate(t *testing.T) {
	expr := NewDivideExpression(NewLiteralExpression(15), NewLiteralExpression(5))
	actual, err := expr.Evaluate()
	if err != nil {
		t.Fatal(err)
	}
	expected := 3
	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestDivideExpressionRange(t *testing.T) {
	type rangeTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected []int
	}
	tests := []rangeTest{
		{
			name:     "TwoInputs",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(15),
			rhs:      NewLiteralExpression(-5),
			expected: []int{-3},
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: []int{0, 1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewDivideExpression(test.lhs, test.rhs)
			actual := GetAllValuesOfRange(expr.Range(), test.name)

			if len(actual) != len(test.expected) {
				t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
			}

			sort.Ints(actual)

			for i := range test.expected {
				if actual[i] != test.expected[i] {
					t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
				}
			}
		})
	}
}

func TestDivideExpressionSimplify(t *testing.T) {
	type simplifyTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected Expression
	}

	tests := []simplifyTest{
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(15),
			rhs:      NewLiteralExpression(5),
			expected: NewLiteralExpression(3),
		},
		{
			name:     "LhsIsZero",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "RhsIsOne",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(1),
			expected: NewInputExpression(0),
		},
		{
			name: "DistributeToAddition",
			lhs:  NewAddExpression(NewInputExpression(0), NewLiteralExpression(15)),
			rhs:  NewLiteralExpression(3),
			// We want (i0/3) + 5
			expected: NewAddExpression(
				NewDivideExpression(
					NewInputExpression(0),
					NewLiteralExpression(3),
				),
				NewLiteralExpression(5),
			),
		},
		{
			name:     "DistributeToAdditionAvoidsIntegerDivisionWeirdness",
			lhs:      NewAddExpression(NewInputExpression(0), NewLiteralExpression(16)),
			rhs:      NewLiteralExpression(3),
			expected: NewDivideExpression(NewAddExpression(NewInputExpression(0), NewLiteralExpression(16)), NewLiteralExpression(3)),
		},
		{
			name:     "DontReduceInputs",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(100),
			expected: NewDivideExpression(NewInputExpression(0), NewLiteralExpression(100)),
		},
		{
			name: "DeepCancellationInMultiplication",
			// (i0 * 10 * i1)
			lhs: NewMultiplyExpression(NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(10)), NewInputExpression(1)),
			// i0
			rhs:      NewInputExpression(0),
			expected: NewMultiplyExpression(NewLiteralExpression(10), NewInputExpression(1)),
		},
		{
			name:     "CancelInputsInMultiplication",
			lhs:      NewInputExpression(0),
			rhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			expected: NewDivideExpression(NewLiteralExpression(1), NewLiteralExpression(20)),
		},
		{
			name:     "CancelLiteralsInMultiplication",
			lhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			rhs:      NewLiteralExpression(20),
			expected: NewInputExpression(0),
		},
		{
			name:     "CancelLiteralsInMultiplicationAvoidsWeirdness",
			lhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			rhs:      NewLiteralExpression(7),
			expected: NewDivideExpression(NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)), NewLiteralExpression(7)),
		},
		{
			name: "DistributeIntoBigGrossThing",
			lhs: NewAddExpression(
				NewInputExpression(0), // Distributing here would potentially lose precision
				NewMultiplyExpression(
					NewEqualsExpression(NewInputExpression(1), NewLiteralExpression(7)),
					NewMultiplyExpression(NewInputExpression(2), NewLiteralExpression(100)),
				),
			),
			rhs: NewLiteralExpression(50),
			expected: NewAddExpression(
				NewMultiplyExpression(
					NewEqualsExpression(NewInputExpression(1), NewLiteralExpression(7)),
					NewMultiplyExpression(NewInputExpression(2), NewLiteralExpression(2)),
				),
				NewDivideExpression(NewInputExpression(0), NewLiteralExpression(50)),
			),
		},
		{
			name: "DistributeIntoAnotherBigGrossThing",
			// (((i0 * 17576) + 17576) + ((i1 * 676) + 7436)) + ((i2 * 26) + 26)
			//
			lhs: NewAddExpression(
				NewAddExpression(
					NewAddExpression(
						NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(17576)),
						NewLiteralExpression(17576),
					),
					NewAddExpression(
						NewMultiplyExpression(NewInputExpression(1), NewLiteralExpression(676)),
						NewLiteralExpression(7436),
					),
				),
				NewAddExpression(
					NewMultiplyExpression(NewInputExpression(2), NewLiteralExpression(26)),
					NewLiteralExpression(26),
				),
			),
			rhs: NewLiteralExpression(26),
			// (i0 * 676) + (i1 * 26) + (i2 + 1) + 962
			expected: NewAddExpression(
				NewAddExpression(
					NewAddExpression(
						NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(676)),
						NewMultiplyExpression(NewInputExpression(1), NewLiteralExpression(26)),
					),
					NewInputExpression(2),
				),
				NewLiteralExpression(963),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewDivideExpression(test.lhs, test.rhs)
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}

}
