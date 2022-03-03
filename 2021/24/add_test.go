package d24

import (
	"testing"
)

func TestNewAddExpression(t *testing.T) {
	expr := NewAddExpression(
		5,
		nil,
		[]*InputExpression{
			NewInputExpression(7).(*InputExpression),
			nil,
			NewInputExpression(9).(*InputExpression),
		},
	)
	expected := "(5 + (i7 + i9))"
	actual := expr.String()
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAddExpressionEvaluate(t *testing.T) {
	expr := NewAddExpression(NewLiteralExpression(5), NewLiteralExpression(10))
	expected := 15
	actual, err := expr.Evaluate()

	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestAddExpressionRange(t *testing.T) {
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
			expected: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(100),
			rhs:      NewLiteralExpression(-500),
			expected: []int{-400},
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(-8),
			expected: []int{-7, -6, -5, -4, -3, -2, -1, 0, 1},
		},
		{
			name:           "TwoValuesWithSameStepOnSharedStepInterval",
			lhs:            NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(26)),
			rhs:            NewMultiplyExpression(NewInputExpression(1), NewLiteralExpression(26)),
			expectedString: "<52..468 step 26>",
		},
		{
			name:           "TwoValuesWithSameStepOnDifferentStepIntervals",
			lhs:            NewAddExpression(NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(26)), NewLiteralExpression(1)),
			rhs:            NewMultiplyExpression(NewInputExpression(1), NewLiteralExpression(26)),
			expectedString: "<53..469 step 26>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewAddExpression(test.lhs, test.rhs)

			if test.expectedString != "" {
				actual := expr.Range().String()
				if actual != test.expectedString {
					t.Errorf("%s: expected %s but got %s", expr.String(), test.expectedString, actual)
				}
			} else {
				actual := GetAllValuesOfRange(expr.Range(), test.name)

				if len(actual) != len(test.expected) {
					t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
				}

				for i := range test.expected {
					if actual[i] != test.expected[i] {
						t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
					}
				}
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
			actual := expr.Simplify([]int{})
			if actual == nil {
				t.Fatal("Simplify() returned nil")
			}
			if test.expected == nil {
				t.Fatal("test.expected is nil")
			}
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}

}
