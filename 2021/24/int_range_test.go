package d24

import "testing"

func TestIntersect(t *testing.T) {
	type intersectTest struct {
		name        string
		lhs         IntRange
		rhs         IntRange
		expected    IntRange
		expectError bool
	}

	tests := []intersectTest{
		{
			name:     "IntersectOnRight",
			lhs:      NewIntRange(0, 5),
			rhs:      NewIntRange(3, 8),
			expected: NewIntRange(3, 5),
		},
		{
			name:     "IntersectOnLeft",
			lhs:      NewIntRange(3, 8),
			rhs:      NewIntRange(0, 5),
			expected: NewIntRange(3, 5),
		},
		{
			name:     "SurroundLeft",
			lhs:      NewIntRange(0, 5),
			rhs:      NewIntRange(1, 2),
			expected: NewIntRange(1, 2),
		},
		{
			name:     "SurroundRight",
			lhs:      NewIntRange(1, 2),
			rhs:      NewIntRange(0, 5),
			expected: NewIntRange(1, 2),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.lhs.Intersect(test.rhs)

			if err == nil {
				if test.expectError {
					t.Errorf("%v ∩ %v should return error, but didn't", test.lhs, test.rhs)
				}
			} else {
				if !test.expectError {
					t.Error(err)
				}
				return
			}

			if actual != test.expected {
				t.Errorf("%v ∩ %v should be %v, but got %v", test.lhs, test.rhs, test.expected, actual)
			}
		})
	}
}

func TestEachWithCustomStep(t *testing.T) {
	r := NewIntRangeWithStep(1, 10, 2)
	actual := make([]int, 0)
	r.Each(func(i int) bool {
		actual = append(actual, i)
		return true
	})

	expected := []int{1, 3, 5, 7, 9}

	if len(actual) != len(expected) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("Expected %v, got %v", expected, actual)
		}
	}

}

func TestSplit(t *testing.T) {
	type splitTest struct {
		name     string
		input    IntRange
		around   int
		expected [3]*IntRange
	}

	r := func(min, max int) *IntRange {
		result := NewIntRange(min, max)
		return &result
	}

	tests := []splitTest{
		{
			name:   "PositiveNegativeAroundZero",
			input:  NewIntRange(-5, 5),
			around: 0,
			expected: [3]*IntRange{
				r(-5, -1),
				r(0, 0),
				r(1, 5),
			},
		},
		{
			name:   "PositiveAroundZero",
			input:  NewIntRange(1, 5),
			around: 0,
			expected: [3]*IntRange{
				nil,
				nil,
				r(1, 5),
			},
		},
		{
			name:   "PositiveIncludingZeroAroundZero",
			input:  NewIntRange(0, 5),
			around: 0,
			expected: [3]*IntRange{
				nil,
				r(0, 0),
				r(1, 5),
			},
		},
		{
			name:   "NegativeAroundZero",
			input:  NewIntRange(-5, -1),
			around: 0,
			expected: [3]*IntRange{
				r(-5, -1),
				nil,
				nil,
			},
		},
		{
			name:   "NegativeIncludingZeroAroundZero",
			input:  NewIntRange(-5, 0),
			around: 0,
			expected: [3]*IntRange{
				r(-5, -1),
				r(0, 0),
				nil,
			},
		},
		{
			name:   "SingleValueRangeAroundItself",
			input:  NewIntRange(5, 5),
			around: 5,
			expected: [3]*IntRange{
				nil,
				r(5, 5),
				nil,
			},
		},
		{
			name:   "SingleValueRangeAroundSmallerValue",
			input:  NewIntRange(5, 5),
			around: 3,
			expected: [3]*IntRange{
				nil,
				nil,
				r(5, 5),
			},
		},
		{
			name:   "SingleValueRangeAroundLargerValue",
			input:  NewIntRange(5, 5),
			around: 8,
			expected: [3]*IntRange{
				r(5, 5),
				nil,
				nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			a, b, c := test.input.Split(test.around)
			actual := []*IntRange{a, b, c}

			for i := 0; i < 3; i++ {
				if test.expected[i] == nil {
					if actual[i] != nil {
						t.Errorf("%d should be nil, but was %v", i, *actual[i])
					}
					continue
				}

				if actual[i] == nil {
					t.Errorf("%d should be %v, but was nil", i, *test.expected[i])
				}

				if *actual[i] != *test.expected[i] {
					t.Errorf("%d should be %v, but was %v", i, *test.expected[i], actual[i])
				}

			}

		})
	}
}
