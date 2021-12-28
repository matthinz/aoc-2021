package d03

import (
	"fmt"
	"strings"
	"testing"
)

func TestCalculateGammaAndEpsilon(t *testing.T) {
	input := strings.Join(
		[]string{
			"00100",
			"11110",
			"10110",
			"10111",
			"10101",
			"01111",
			"00111",
			"11100",
			"10000",
			"11001",
			"00010",
			"01010",
		},
		"\n",
	)

	numbers := parseInput(strings.NewReader(input))

	gamma, epsilon := calculateGammaAndEpsilon(numbers)

	expectedGamma := 22
	if gamma != expectedGamma {
		t.Error(fmt.Sprintf("Expected gamma %d (%b), got %d (%b)", expectedGamma, expectedGamma, gamma, gamma))
	}

	expectedEpsilon := 9
	if epsilon != expectedEpsilon {
		t.Error(fmt.Sprintf("Expected epsilon %d (%b), got %d (%b)", expectedEpsilon, expectedEpsilon, epsilon, epsilon))
	}
}
