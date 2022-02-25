package d24

import (
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

func TestEqualsExpressionRange(t *testing.T) {
	expr := NewEqualsExpression(NewLiteralExpression(5), NewLiteralExpression(8))
	expected := []int{0, 1}
	actual := GetAllValuesOfRange(expr.Range(), "TestEqualsExpressionRange")

	if len(actual) != len(expected) {
		t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), expected, len(expected), actual, len(actual))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), expected, len(expected), actual, len(actual))
		}
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
