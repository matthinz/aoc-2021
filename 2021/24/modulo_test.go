package d24

import (
	"sort"
	"testing"
)

func TestModuloExpressionEvaluate(t *testing.T) {
	type evaluateTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		inputs   []int
		expected int
	}

	tests := []evaluateTest{
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(15),
			rhs:      NewLiteralExpression(4),
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)
			actual, err := expr.Evaluate()
			if err != nil {
				t.Fatal(err)
			}
			if actual != test.expected {
				t.Errorf("%s: for inputs %v expected %d but got %d", expr.String(), test.inputs, test.expected, actual)
			}
		})
	}
}

func TestModuloExpressionRange(t *testing.T) {
	type rangeTest struct {
		name           string
		lhs            Expression
		rhs            Expression
		expected       []int
		expectedString string
	}
	tests := []rangeTest{
		{
			name:     "TwoInputs",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(3),
			expected: []int{2},
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: []int{0, 1, 2},
		},
		{
			name:     "NegativeLiteralLhsRhsInput",
			lhs:      NewLiteralExpression(-5),
			rhs:      NewInputExpression(0),
			expected: []int{-5, -2, -1, 0},
		},
		{
			name:     "LhsInputNegativeLiteralRhs",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(-5),
			expected: []int{0, 1, 2, 3, 4},
		},
		{
			name:     "LargeLhsLargeRhs",
			lhs:      NewLiteralExpression(1234567890),
			rhs:      NewMultiplyExpression(NewLiteralExpression(251), NewInputExpression(0)),
			expected: []int{43, 294, 545, 1298, 1800},
		},
		{
			name:     "NegativeInputLhsRhsLiteral",
			lhs:      NewMultiplyExpression(NewLiteralExpression(-2), NewInputExpression(0)),
			rhs:      NewLiteralExpression(4),
			expected: []int{-2, 0},
		},
		{
			name:           "LargeLhsSmallRhs",
			lhs:            NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(100000)),
			rhs:            NewLiteralExpression(26),
			expected:       []int{2, 4, 6, 8, 10, 12, 16, 20, 24},
			expectedString: "<<2..12 step 2>,<16..24 step 4>>",
		},
		{
			name:     "RhsEqualToLhsStep",
			lhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(5)),
			rhs:      NewLiteralExpression(5),
			expected: []int{0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)

			if test.expectedString == "" {

				actual := GetAllValuesOfRange(expr.Range(), test.name)

				if len(actual) != len(test.expected) {
					t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
				}

				// NOTE: The order in which we get values is not stable, so we have to sort to compare
				sort.Ints(actual)

				for i := range test.expected {
					if actual[i] != test.expected[i] {
						t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
					}
				}
			} else {
				actual := expr.Range().String()
				if actual != test.expectedString {
					t.Errorf("%s: expected %s but got %s", expr.String(), test.expectedString, actual)
				}
			}

		})
	}
}

func TestModuloExpressionSimplify(t *testing.T) {
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
			rhs:      NewLiteralExpression(3),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(8),
			expected: NewModuloExpression(NewInputExpression(0), NewLiteralExpression(8)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)
			actual := expr.Simplify([]int{})
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}
