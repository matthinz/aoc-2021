package d24

import (
	"fmt"
	"testing"
)

func TestMergeContinuousRanges(t *testing.T) {

	type test struct {
		ranges   []ContinuousRange
		expected string
	}

	tests := []test{
		{
			ranges: []ContinuousRange{
				newContinuousRange(0, 5, 1),
				newContinuousRange(4, 10, 2),
			},
			expected: "<<0..6>,<8,10>>",
		},
		{
			ranges: []ContinuousRange{
				newContinuousRange(1, 9, 1),
				newContinuousRange(2, 18, 2),
				newContinuousRange(3, 27, 3),
				newContinuousRange(4, 36, 4),
			},
			expected: "<<1..10>,<12..20 step 2>,<15..27 step 3>,<28,32,36>>",
		},
		{
			ranges: []ContinuousRange{
				newContinuousRange(1665, 1673, 1),
				newContinuousRange(1691, 1699, 1),
				newContinuousRange(1717, 1725, 1),
			},
			expected: "<<1665..1673>,<1691..1699>,<1717..1725>>",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := NewContinuousRangeSet(test.ranges)
			actual := r.String()
			if actual != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, actual)
			}
		})
	}

}
