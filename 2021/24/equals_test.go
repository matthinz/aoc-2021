package d24

import (
	"testing"
)

func TestEqualsExpressionEvaluate(t *testing.T) {
	type evaluateTest struct {
		name     string
		lhs      Expression
		rhs      Expression
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewEqualsExpression(test.lhs, test.rhs)
			actual, err := expr.Evaluate()
			if err != nil {
				t.Fatal(err)
			}
			if actual != test.expected {
				t.Errorf("%s: expected %d but got %d", expr.String(), test.expected, actual)
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
		// TODO: Re-enable this (should be able to quickly check for nonintersection)
		// {
		// 	name:     "InputAndLiteralOutsideRange",
		// 	lhs:      NewInputExpression(0),
		// 	rhs:      NewLiteralExpression(500),
		// 	expected: NewLiteralExpression(0),
		// },
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
			actual := expr.Simplify(map[int]int{})
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}
