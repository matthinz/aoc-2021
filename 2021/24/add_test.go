package d24

import (
	"log"
	"testing"
)

func TestAddExpressionEvaluate(t *testing.T) {
	expr := NewAddExpression(NewLiteralExpression(5), NewInputExpression(0))
	expected := 15
	actual := expr.Evaluate([]int{10})
	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestAddExpressionFindInputs(t *testing.T) {
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
			target:   12,
			decider:  PreferFirstSetOfInputs,
			expected: []int{7},
		},
		{
			name:     "LhsInputRhsLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(5),
			target:   12,
			decider:  PreferFirstSetOfInputs,
			expected: []int{7},
		},
		{
			name:     "TwoInputsThatMustBeEqual",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(0),
			target:   12,
			decider:  PreferFirstSetOfInputs,
			expected: []int{6},
		},
		{
			name:     "TwoInputsThatMakeLargestNumber",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			target:   12,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9, 3},
		},
		{
			name:        "TwoInputsThatCantAddUp",
			lhs:         NewInputExpression(0),
			rhs:         NewInputExpression(1),
			target:      100,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewAddExpression(test.lhs, test.rhs)
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

func TestAddExpressionRange(t *testing.T) {
	type rangeTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected IntRange
	}
	tests := []rangeTest{
		{
			name:     "TwoInputs",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			expected: NewIntRange(2, 18),
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(100),
			rhs:      NewLiteralExpression(-500),
			expected: NewIntRange(-400, -400),
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(-8),
			expected: NewIntRange(-7, 1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewAddExpression(test.lhs, test.rhs)
			actual := expr.Range()
			if actual != test.expected {
				t.Errorf("%s: expected range %v but got %v", expr.String(), test.expected, actual)
			}
		})
	}
}

func TestAddExpressionSimplify(t *testing.T) {
	type simplifyTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected Expression
	}

	tests := []simplifyTest{
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(8),
			expected: NewLiteralExpression(13),
		},
		{
			name:     "LhsIsZero",
			lhs:      NewLiteralExpression(0),
			rhs:      NewLiteralExpression(8),
			expected: NewLiteralExpression(8),
		},
		{
			name:     "RhsIsZero",
			lhs:      NewLiteralExpression(8),
			rhs:      NewLiteralExpression(0),
			expected: NewLiteralExpression(8),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewAddExpression(test.lhs, test.rhs)
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}

}
