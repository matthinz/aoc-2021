package d24

import (
	"fmt"
	"testing"
)

func TestContinuousRangeIntersect(t *testing.T) {
	type test struct {
		name     string
		lhs      ContinuousRange
		rhs      ContinuousRange
		expected bool
	}

	tests := []test{
		{
			lhs:      newContinuousRange(8, 8, 1),
			rhs:      newContinuousRange(1, 9, 1),
			expected: true,
		},
		{
			lhs:      newContinuousRange(1, 9, 1),
			rhs:      newContinuousRange(8, 8, 1),
			expected: true,
		},
		{
			lhs:      newContinuousRange(0, 0, 1),
			rhs:      newContinuousRange(1, 9, 1),
			expected: false,
		},
		{
			lhs:      newContinuousRange(0, 10, 2),
			rhs:      newContinuousRange(1, 11, 2),
			expected: false,
		},
		{
			lhs:      newContinuousRange(0, 5, 1),
			rhs:      newContinuousRange(4, 10, 1),
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s,%s", test.lhs, test.rhs), func(t *testing.T) {

			actual := test.lhs.(*continuousRange).Intersects(test.rhs.(*continuousRange))
			if actual && !test.expected {
				t.Errorf("%s should NOT intersect %s", test.lhs, test.rhs)
			} else if test.expected && !actual {
				t.Errorf("%s should intersect %s", test.lhs, test.rhs)
			}
		})
	}

}
