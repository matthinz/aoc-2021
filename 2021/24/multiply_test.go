package d24

import "testing"

func TestMultiplyExpressionEvaluate(t *testing.T) {
	expr := NewMultiplyExpression(NewLiteralExpression(15), NewInputExpression(0))
	expected := 45
	actual := expr.Evaluate([]int{3})
	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestMultiplyExpressionFindInputs(t *testing.T) {
	t.Skip()
}

func TestMultiplyExpressionRange(t *testing.T) {
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
			expected: NewIntRange(1, 81),
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(5),
			rhs:      NewLiteralExpression(3),
			expected: NewIntRange(15, 15),
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: NewIntRange(3, 27),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewMultiplyExpression(test.lhs, test.rhs)
			actual := expr.Range()
			if actual != test.expected {
				t.Errorf("%s: expected range %v but got %v", expr.String(), test.expected, actual)
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
