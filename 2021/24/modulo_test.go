package d24

import "testing"

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
	t.Skip()
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
