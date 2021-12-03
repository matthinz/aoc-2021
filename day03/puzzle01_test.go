package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCalculateGammaAndEpsilon(t *testing.T) {
	input := []string{
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
	}

	numbers := parseInput(input)

	gamma, epsilon := CalculateGammaAndEpsilon(numbers)

	expectedGamma := 22
	if gamma != expectedGamma {
		t.Error(fmt.Sprintf("Expected gamma %d (%b), got %d (%b)", expectedGamma, expectedGamma, gamma, gamma))
	}

	expectedEpsilon := 9
	if epsilon != expectedEpsilon {
		t.Error(fmt.Sprintf("Expected epsilon %d (%b), got %d (%b)", expectedEpsilon, expectedEpsilon, epsilon, epsilon))
	}
}

func parseInput(binaryNumbers []string) []int {
	var result []int

	for _, token := range binaryNumbers {
		value, err := strconv.ParseInt(token, 2, 64)
		if err != nil {
			panic(err)
		}
		result = append(result, int(value))
	}

	return result
}
