package d24

import (
	"log"
	"sort"
	"testing"
)

func TestDivideExpressionEvaluate(t *testing.T) {
	expr := NewDivideExpression(NewLiteralExpression(15), NewInputExpression(0))
	expected := 5
	actual := expr.Evaluate([]int{3})
	if actual != expected {
		t.Errorf("%s: expected %d, got %d", expr.String(), expected, actual)
	}
}

func TestDivideExpressionFindInputs(t *testing.T) {
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
			lhs:      NewLiteralExpression(12),
			rhs:      NewInputExpression(0),
			target:   4,
			decider:  PreferFirstSetOfInputs,
			expected: []int{3},
		},
		{
			name:     "LhsInputRhsLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			target:   3,
			decider:  PreferFirstSetOfInputs,
			expected: []int{9},
		},
		{
			name:     "TwoInputsThatMustBeEqual",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(0),
			target:   1,
			decider:  PreferFirstSetOfInputs,
			expected: []int{1},
		},
		{
			name:     "TwoInputsThatMakeLargestNumber",
			lhs:      NewInputExpression(0),
			rhs:      NewInputExpression(1),
			target:   9,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9, 1},
		},
		{
			name:        "DivideByZero",
			lhs:         NewInputExpression(0),
			rhs:         NewLiteralExpression(0),
			target:      4,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
		{
			name:     "ZeroNumerator",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			target:   0,
			decider:  PreferFirstSetOfInputs,
			expected: []int{1},
		},
		{
			name:     "ZeroNumeratorLargestInput",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			target:   0,
			decider:  PreferInputsThatMakeLargerNumber,
			expected: []int{9},
		},
		{
			name:     "DivideLiteralBy1",
			lhs:      NewLiteralExpression(3),
			rhs:      NewLiteralExpression(1),
			target:   3,
			decider:  PreferFirstSetOfInputs,
			expected: []int{},
		},
		{
			name:        "InvalidLiteralDivide",
			lhs:         NewLiteralExpression(6),
			rhs:         NewLiteralExpression(3),
			target:      12,
			decider:     PreferFirstSetOfInputs,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewDivideExpression(test.lhs, test.rhs)
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

func TestDivideExpressionRange(t *testing.T) {
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
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:     "TwoLiterals",
			lhs:      NewLiteralExpression(15),
			rhs:      NewLiteralExpression(-5),
			expected: []int{-3},
		},
		{
			name:     "InputAndLiteral",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(3),
			expected: []int{0, 1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewDivideExpression(test.lhs, test.rhs)
			actual := GetAllValuesOfRange(expr.Range())

			if len(actual) != len(test.expected) {
				t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
			}

			sort.Ints(actual)

			for i := range test.expected {
				if actual[i] != test.expected[i] {
					t.Fatalf("%s: expected range %v (%d) but got %v (%d)", expr.String(), test.expected, len(test.expected), actual, len(actual))
				}
			}
		})
	}
}

func TestDivideExpressionSimplify(t *testing.T) {
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
			rhs:      NewLiteralExpression(5),
			expected: NewLiteralExpression(3),
		},
		{
			name:     "LhsIsZero",
			lhs:      NewLiteralExpression(0),
			rhs:      NewInputExpression(0),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "RhsIsOne",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(1),
			expected: NewInputExpression(0),
		},
		{
			name:     "DistributeToAddition",
			lhs:      NewAddExpression(NewInputExpression(0), NewLiteralExpression(15)),
			rhs:      NewLiteralExpression(3),
			expected: NewAddExpression(NewDivideExpression(NewInputExpression(0), NewLiteralExpression(3)), NewLiteralExpression(5)),
		},
		{
			name:     "DistributeToAdditionAvoidsIntegerDivisionWeirdness",
			lhs:      NewAddExpression(NewInputExpression(0), NewLiteralExpression(16)),
			rhs:      NewLiteralExpression(3),
			expected: NewDivideExpression(NewAddExpression(NewInputExpression(0), NewLiteralExpression(16)), NewLiteralExpression(3)),
		},
		{
			name:     "LargeRhsReducesToZero",
			lhs:      NewInputExpression(0),
			rhs:      NewLiteralExpression(100),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "LargeRhsRangeReducesToZero",
			lhs:      NewInputExpression(0),
			rhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			expected: NewLiteralExpression(0),
		},
		{
			name:     "CancelElementsInMultiplication",
			lhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			rhs:      NewLiteralExpression(20),
			expected: NewInputExpression(0),
		},
		{
			name:     "CancelElementsInMultiplicationAvoidsWeirdness",
			lhs:      NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)),
			rhs:      NewLiteralExpression(7),
			expected: NewDivideExpression(NewMultiplyExpression(NewInputExpression(0), NewLiteralExpression(20)), NewLiteralExpression(7)),
		},
		{
			name: "DistributeIntoBigGrossThing",
			lhs: NewAddExpression(
				NewInputExpression(0),
				NewMultiplyExpression(
					NewEqualsExpression(NewInputExpression(1), NewLiteralExpression(7)),
					NewMultiplyExpression(NewInputExpression(2), NewLiteralExpression(100)),
				),
			),
			rhs: NewLiteralExpression(50),
			expected: NewMultiplyExpression(
				NewEqualsExpression(NewInputExpression(1), NewLiteralExpression(7)),
				NewMultiplyExpression(NewInputExpression(2), NewLiteralExpression(2)),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := NewDivideExpression(test.lhs, test.rhs)
			actual := expr.Simplify()
			if actual.String() != test.expected.String() {
				t.Errorf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}

}
