package d24

import (
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
		{
			name:     "LiteralAndInput",
			lhs:      NewLiteralExpression(15),
			rhs:      NewInputExpression(0),
			inputs:   []int{4},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)
			actual := expr.Evaluate(test.inputs)
			if actual != test.expected {
				t.Errorf("%s: for inputs %v expected %d but got %d", expr.String(), test.inputs, test.expected, actual)
			}
		})
	}
}

func TestModuloExpressionFindInputs(t *testing.T) {
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
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(10),
			rhs:      NewLiteralExpression(3),
			target:   1,
			decider:  PreferFirstSetOfInputs,
			expected: []int{},
		},
		{
			name:        "TwoLiteralsCantHitTarget",
			lhs:         NewLiteralExpression(10),
			rhs:         NewLiteralExpression(3),
			target:      45,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
		{
			name:     "LhsLiteralRhsInput",
			lhs:      NewLiteralExpression(10),
			rhs:      NewInputExpression(0),
			target:   1,
			decider:  PreferFirstSetOfInputs,
			expected: []int{3},
		},
		{
			name:     "LhsInputRhsLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(6),
			target:   1,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{7},
		},
		{
			name:     "ComplexWithBigNumbers",
			lhs:      NewLiteralExpression(12345678),
			rhs:      NewAddExpression(NewInputExpression(0), NewLiteralExpression(25)),
			target:   15,
			decider:  PreferFirstSetOfInputs,
			expected: []int{8},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)
			actualMap, err := expr.FindInputs(test.target, test.decider)

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

func TestModuloExpressionRange(t *testing.T) {
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
			expected: NewIntRange(0, 8),
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(3),
			expected: NewIntRange(2, 2),
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: NewIntRange(0, 2),
		},
		{
			name:     "NegativeLiteralLhsRhsInput",
			lhs:      NewLiteralExpression(-5),
			rhs:      NewInputExpression(0),
			expected: NewIntRange(-5, 0),
		},
		{
			name:     "LhsInputNegativeLiteralRhs",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(-5),
			expected: NewIntRange(0, 4),
		},
		{
			name:     "LargeLhsLargeRhs",
			lhs:      NewLiteralExpression(1234567890),
			rhs:      NewMultiplyExpression(NewLiteralExpression(251), NewInputExpression(0)),
			expected: NewIntRange(43, 1800),
		},
		{
			name:     "NegativeInputLhsRhsLiteral",
			lhs:      NewMultiplyExpression(NewLiteralExpression(-2), NewInputExpression(0)),
			rhs:      NewLiteralExpression(4),
			expected: NewIntRange(-2, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewModuloExpression(test.lhs, test.rhs)
			actual := expr.Range()

			if actual != test.expected {
				t.Errorf("%s: expected range %v but got %v", expr.String(), test.expected, actual)
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
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}
