package d24

import (
	"log"
	"testing"
)

func TestMultiplyExpressionEvaluate(t *testing.T) {
	expr := NewMultiplyExpression(NewLiteralExpression(15), NewInputExpression(0))
	expected := 45
	actual := expr.Evaluate([]int{3})
	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestMultiplyExpressionFindInputs(t *testing.T) {
	type findInputsTest struct {
		name        string
		lhs         Expression
		rhs         Expression
		target      int
		decider     InputDecider
		expected    []int
		expectError bool
	}

	tests := []findInputsTest{
		{
			name:     "LhsLiteralRhsInput",
			lhs:      NewLiteralExpression(5),
			rhs:      NewInputExpression(0),
			target:   15,
			decider:  PreferFirstSetOfInputs,
			expected: []int{3},
		},
		{
			name:     "LhsInputRhsLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(5),
			target:   15,
			decider:  PreferFirstSetOfInputs,
			expected: []int{3},
		},
		{
			name:     "TwoInputsThatMustBeEqual",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(0),
			target:   16,
			decider:  PreferFirstSetOfInputs,
			expected: []int{4},
		},
		{
			name:     "TwoInputsThatMakeLargestNumber",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			target:   12,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{6, 2},
		},
		{
			name:     "LhsInputToZero",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(0),
			target:   0,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9},
		},
		{
			name:     "RhsInputToZero",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			target:   0,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewMultiplyExpression(test.lhs, test.rhs)
			actualMap, err := expr.FindInputs(test.target, test.decider, log.Default())

			if test.expectError && err == nil {
				t.Fatalf("Expected test to error but it didn't")
			} else if err != nil && !test.expectError {
				t.Fatal(err)
			}

			if test.expectError {
				return
			}

			actual := make([]int, len(actualMap))
			for index, value := range actualMap {
				actual[index] = value
			}

			if len(actual) != len(test.expected) {
				t.Fatalf("Wrong # of items in result. Expected %v (%d), got %v (%d)", test.expected, len(test.expected), actual, len(actual))
			}

			for i := range test.expected {
				if actual[i] != test.expected[i] {
					t.Errorf("Item %d is wrong. Expected %d, got %d", i, test.expected[i], actual[i])
				}
			}
		})
	}
}

func TestMultiplyExpressionRange(t *testing.T) {
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
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 2, 4, 6, 8, 10, 12, 14, 16, 18, 3, 6, 9, 12, 15, 21, 24, 27, 4, 8, 12, 16, 20, 24, 28, 32, 36, 5, 10, 15, 20, 25, 30, 35, 40, 45, 6, 12, 18, 24, 30, 36, 42, 48, 54, 7, 14, 21, 28, 35, 42, 49, 56, 63, 8, 16, 24, 32, 40, 48, 56, 64, 72, 9, 18, 27, 36, 45, 54, 63, 81},
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(3),
			expected: []int{15},
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: []int{3, 6, 9, 12, 15, 18, 21, 24, 27},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewMultiplyExpression(test.lhs, test.rhs)
			actual := GetAllValuesOfRange(expr.Range())

			if len(actual) != len(test.expected) {
				t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
			}

			for i := range test.expected {
				if actual[i] != test.expected[i] {
					t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
				}
			}
		})
	}
}

func TestMultiplyExpressionSimplify(t *testing.T) {
	type simplifyTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected Expression
	}

	tests := []simplifyTest{
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(3),
			rhs:      NewLiteralExpression(5),
			expected: NewLiteralExpression(15),
		},
		{
			name:     "LhsIsZero",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "RLhsIsZero",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(0),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "LhsIsOne",
			lhs:      NewLiteralExpression(1),
			rhs:      NewInputExpression(0),
			expected: NewInputExpression(0),
		},
		{
			name:     "RhsIsOne",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(1),
			expected: NewInputExpression(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewMultiplyExpression(test.lhs, test.rhs)
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}
