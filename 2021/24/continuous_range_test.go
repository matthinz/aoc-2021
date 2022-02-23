package d24

import "testing"

func TestSplitContinuousRange(t *testing.T) {

	type continuousRangeSplitTest struct {
		name     string
		min, max int
		step     int
		around   Range
		expected []*continuousRange
	}

	tests := []continuousRangeSplitTest{
		{
			name:   "AroundNumberNotInRangeDueToStep",
			min:    2,
			max:    10,
			step:   2,
			around: newContinuousRange(5, 5, 1),
			expected: []*continuousRange{
				newContinuousRange(2, 4, 2),
				nil,
				newContinuousRange(6, 10, 2),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := newContinuousRange(test.min, test.max, test.step)
			before, at, after := r.Split(test.around)

			if test.expected[0] == nil {
				if before != nil {
					t.Errorf("Expected before to be nil, but got %v", before)
				}
			} else {
				if before == nil {
					t.Errorf("Expected before to be %v, but got nil", *test.expected[0])
				} else {
					if !RangesEqual(before, test.expected[0]) {
						t.Errorf("Expected before to be %v, but got %v", *test.expected[0], before)
					}
				}
			}

			if test.expected[1] == nil {
				if at != nil {
					t.Errorf("Expected at to be nil, but got %v", at)
				}
			} else {
				if at == nil {
					t.Errorf("Expected at to be %v, but got nil", *test.expected[1])
				} else {
					if !RangesEqual(at, test.expected[1]) {
						t.Errorf("Expected at to be %v, but got %v", *test.expected[1], at)
					}
				}
			}

			if test.expected[2] == nil {
				if after != nil {
					t.Errorf("Expected after to be nil, but got %v", after)
				}
			} else {
				if after == nil {
					t.Errorf("Expected after to be %v, but got nil", *test.expected[2])
				} else {
					if !RangesEqual(after, test.expected[2]) {
						t.Errorf("Expected after to be %v, but got %v", *test.expected[2], after)
					}
				}
			}

		})
	}

}
