package main

import (
	"fmt"
	"testing"
)

func TestStep(t *testing.T) {
	steps := [][][]int{
		[][]int{
			[]int{5, 4, 8, 3, 1, 4, 3, 2, 2, 3},
			[]int{2, 7, 4, 5, 8, 5, 4, 7, 1, 1},
			[]int{5, 2, 6, 4, 5, 5, 6, 1, 7, 3},
			[]int{6, 1, 4, 1, 3, 3, 6, 1, 4, 6},
			[]int{6, 3, 5, 7, 3, 8, 5, 4, 7, 8},
			[]int{4, 1, 6, 7, 5, 2, 4, 6, 4, 5},
			[]int{2, 1, 7, 6, 8, 4, 1, 7, 2, 1},
			[]int{6, 8, 8, 2, 8, 8, 1, 1, 3, 4},
			[]int{4, 8, 4, 6, 8, 4, 8, 5, 5, 4},
			[]int{5, 2, 8, 3, 7, 5, 1, 5, 2, 6},
		},
		[][]int{
			[]int{6, 5, 9, 4, 2, 5, 4, 3, 3, 4},
			[]int{3, 8, 5, 6, 9, 6, 5, 8, 2, 2},
			[]int{6, 3, 7, 5, 6, 6, 7, 2, 8, 4},
			[]int{7, 2, 5, 2, 4, 4, 7, 2, 5, 7},
			[]int{7, 4, 6, 8, 4, 9, 6, 5, 8, 9},
			[]int{5, 2, 7, 8, 6, 3, 5, 7, 5, 6},
			[]int{3, 2, 8, 7, 9, 5, 2, 8, 3, 2},
			[]int{7, 9, 9, 3, 9, 9, 2, 2, 4, 5},
			[]int{5, 9, 5, 7, 9, 5, 9, 6, 6, 5},
			[]int{6, 3, 9, 4, 8, 6, 2, 6, 3, 7},
		},
		[][]int{
			[]int{8, 8, 0, 7, 4, 7, 6, 5, 5, 5},
			[]int{5, 0, 8, 9, 0, 8, 7, 0, 5, 4},
			[]int{8, 5, 9, 7, 8, 8, 9, 6, 0, 8},
			[]int{8, 4, 8, 5, 7, 6, 9, 6, 0, 0},
			[]int{8, 7, 0, 0, 9, 0, 8, 8, 0, 0},
			[]int{6, 6, 0, 0, 0, 8, 8, 9, 8, 9},
			[]int{6, 8, 0, 0, 0, 0, 5, 9, 4, 3},
			[]int{0, 0, 0, 0, 0, 0, 7, 4, 5, 6},
			[]int{9, 0, 0, 0, 0, 0, 0, 8, 7, 6},
			[]int{8, 7, 0, 0, 0, 0, 6, 8, 4, 8},
		},
	}

	for i := 0; i < len(steps)-1; i++ {
		input := steps[i]
		expected := steps[i+1]

		actual, _ := step(input)

		t.Run(fmt.Sprintf("step %d", i), func(t *testing.T) {
			assertEq(t, expected, actual)
		})
	}

}

func assertEq(t *testing.T, expected, actual [][]int) {

	if len(expected) != len(actual) {
		t.Fatalf("Different heights: %d vs %d", len(expected), len(actual))
	}

	ok := true

	for y := 0; y < len(expected); y++ {

		expectedRow := expected[y]
		actualRow := actual[y]

		if len(expectedRow) != len(actualRow) {
			t.Fatalf("Row %d: expected length %d, got length %d", y, len(expectedRow), len(actualRow))
		}

		for x := 0; x < len(expectedRow); x++ {

			expectedOctopus := expectedRow[x]
			actualOctopus := actualRow[x]

			if actualOctopus != expectedOctopus {
				t.Errorf("%d,%d: expected %d, got %d", x, y, expectedOctopus, actualOctopus)
				ok = false
			}
		}
	}

	if !ok {
		fmt.Println("EXPECTED")
		printGrid(expected)
		fmt.Println("ACTUAL")
		printGrid(actual)
	}

}
