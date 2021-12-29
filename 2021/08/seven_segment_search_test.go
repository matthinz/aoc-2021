package d08

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := strings.Split("acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab", " ")

	expectedNumbers :=
		[]int{
			8,
			5,
			2,
			3,
			7,
			9,
			6,
			4,
			0,
			1,
		}

	actual := interpretSignalValues(input)

	if len(actual) != 10 {
		t.Fatalf("wrong # of actual things: %d", len(actual))
	}

	ok := true
	for i := range expectedNumbers {
		if expectedNumbers[i] != actual[i] {
			t.Logf("%d: expected %d, got %d", i, expectedNumbers[i], actual[i])
			ok = false
		}
	}

	if !ok {
		t.Fatal("failed")
	}
}
