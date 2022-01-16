package d24

import (
	"log"
	"testing"
)

func TestEqualsExpressionEvaluate(t *testing.T) {
	type evaluateTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		inputs   []int
		expected int
	}

	tests := []evaluateTest{
		{
			name:     "EqualLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(5),
			expected: 1,
		},
		{
			name:     "UnequalLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(8),
			expected: 0,
		},
		{
			name:     "EqualLiteralAndInput",
			lhs:      NewLiteralExpression(5),
			rhs:      NewInputExpression(0),
			inputs:   []int{5},
			expected: 1,
		},
		{
			name:     "UnequalLiteralAndInput",
			lhs:      NewLiteralExpression(5),
			rhs:      NewInputExpression(0),
			inputs:   []int{3},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewEqualsExpression(test.lhs, test.rhs)
			actual := expr.Evaluate(test.inputs)
			if actual != test.expected {
				t.Errorf("%s: for inputs %v expected %d but got %d", expr.String(), test.inputs, test.expected, actual)
			}
		})
	}
}

func TestEqualsExpressionFindInputs(t *testing.T) {
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
			name:     "TwoEqualLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(5),
			target:   1,
			decider:  PreferFirstSetOfInputs,
			expected: []int{},
		},
		{
			name:     "TwoUnequalLiterals",
			lhs:      NewLiteralExpression(3),
			rhs:      NewLiteralExpression(5),
			target:   0,
			decider:  PreferFirstSetOfInputs,
			expected: []int{},
		},
		{
			name:        "TwoEqualLiteralsSeekingNotEqual",
			lhs:         NewLiteralExpression(5),
			rhs:         NewLiteralExpression(5),
			target:      0,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
		{
			name:        "TwoNonEqualLiteralsSeekingEqual",
			lhs:         NewLiteralExpression(3),
			rhs:         NewLiteralExpression(5),
			target:      1,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
		{
			name:     "TwoDifferentInputsSeekingEquality",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			target:   1,
			decider:  PreferFirstSetOfInputs,
			expected: []int{1, 1},
		},
		{
			name:     "TwoDifferentInputsSeekingInequality",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			target:   0,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9, 8},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewEqualsExpression(test.lhs, test.rhs)
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

			t.Logf("%s evaluates to %d for %v (target was %d)", expr.String(), expr.Evaluate(actual), actual, test.target)

			loggedInputs := false
			for i := range test.expected {
				if !loggedInputs {
					t.Logf("Expected: %v", test.expected)
					t.Logf("Actual: %v", actual)
					loggedInputs = true
				}
				if actual[i] != test.expected[i] {
					t.Errorf("Item %d is wrong. Expected %d, got %d", i, test.expected[i], actual[i])
				}
			}
		})
	}
}

func TestEqualsExpressionRange(t *testing.T) {
	expr := NewEqualsExpression(NewLiteralExpression(5), NewLiteralExpression(8))
	expected := NewIntRange(0, 1)
	actual := expr.Range()
	if actual != expected {
		t.Errorf("Expected range %v, but got %v", expected, actual)
	}
}

func TestEqualsExpressionSimplify(t *testing.T) {
	type simplifyTest struct {
		name     string
		lhs      Expression
		rhs      Expression
		expected Expression
	}

	tests := []simplifyTest{
		{
			name:     "TwoEqualLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(5),
			expected: NewLiteralExpression(1),
		},
		{
			name:     "TwoUnequalLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(8),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "InputAndLiteralOutsideRange",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(500),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "InputAndLiteralInsideRange",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(8),
			expected: NewEqualsExpression(NewInputExpression(0), NewLiteralExpression(8)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewEqualsExpression(test.lhs, test.rhs)
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}
