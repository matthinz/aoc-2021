package d24

import "testing"

func TestNewCompoundRangeFromContinuousRanges(t *testing.T) {
	type test struct {
		name     string
		a        Range
		b        Range
		expected string
	}

	tests := []test{
		{
			name:     "SameStepAdjacent",
			a:        newContinuousRange(0, 1, 1),
			b:        newContinuousRange(1, 9, 1),
			expected: "<0..9>",
		},
		{
			name:     "SameStepIntersectLeft",
			a:        newContinuousRange(0, 2, 1),
			b:        newContinuousRange(1, 9, 1),
			expected: "<0..9>",
		},
		{
			name:     "SameStepIntersectRight",
			a:        newContinuousRange(8, 12, 1),
			b:        newContinuousRange(1, 9, 1),
			expected: "<1..12>",
		},
		{
			name:     "DifferentStepIntersect",
			a:        newContinuousRange(2, 8, 2),
			b:        newContinuousRange(3, 9, 3),
			expected: "<<2..8 step 2>,<3,6,9>>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := newCompoundRange(test.a, test.b)
			if actual.String() != test.expected {
				t.Errorf("%s + %s: expected '%s' but got '%s'", test.a, test.b, test.expected, actual.String())
			}
		})
	}

}
