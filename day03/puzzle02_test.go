package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParsingWorksHowIThink(t *testing.T) {
	value, err := strconv.ParseInt("01010", 2, 64)
	if err != nil {
		t.Error(err)
	}
	if value != 10 {
		t.Error(fmt.Sprintf("Expected 10, got %d", value))
	}
}

func TestFindO2GeneratorRating(t *testing.T) {
	input := parseInput([]string{
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
	})

	actual := FindO2GeneratorRating(input)
	expected := 23
	if actual != expected {
		t.Error(fmt.Sprintf("Got the wrong rating! Expected %d (%b), got %d (%b)", expected, expected, actual, actual))
	}
}

func TestFindCo2ScrubberRating(t *testing.T) {
	input := parseInput([]string{
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
	})

	actual := FindCo2ScrubberRating(input)
	expected := 10
	if actual != expected {
		t.Error(fmt.Sprintf("Got the wrong rating! Expected %d (%b), got %d (%b)", expected, expected, actual, actual))
	}
}

func TestFindMostAndLeastCommonBits(t *testing.T) {
	input := parseInput([]string{
		"1001",
		"1010",
		"0010",
		"1111",
	})

	expectedMostCommon := []int{
		1,
		1,
		0,
		1,
	}

	expectedLeastCommon := []int{
		0,
		0,
		1,
		0,
	}

	for i := 0; i < 4; i++ {
		mostCommon, leastCommon := FindMostAndLeastCommonBits(input, i)

		if mostCommon != expectedMostCommon[i] {
			t.Error(fmt.Sprintf("mostCommon should be %d at position %d", expectedMostCommon[i], i))
		}
		if leastCommon != expectedLeastCommon[i] {
			t.Error(fmt.Sprintf("leastCommon should be %d at position %d", expectedLeastCommon[i], i))
		}
	}
}

func TestFilter(t *testing.T) {
	input := parseInput([]string{
		"1001",
		"1010",
		"0010",
		"1111",
	})

	actual := FilterNumbersByBitAtPosition(input, 1, 0)
	expected := parseInput([]string{
		"1001",
		"1111",
	})

	assertEqual(t, actual, expected)

	actual = FilterNumbersByBitAtPosition(input, 0, 0)
	expected = parseInput([]string{
		"1010",
		"0010",
	})
	assertEqual(t, actual, expected)

	actual = FilterNumbersByBitAtPosition(input, 0, 3)
	expected = parseInput([]string{
		"1001",
		"1010",
		"0010",
	})
}

func assertEqual(t *testing.T, actual []int, expected []int) {

	if len(actual) != len(expected) {
		t.Error(actual)
		return
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Error(fmt.Sprintf("%d: expected %d (%b), got %d (%b)", i, expected[i], expected[i], actual[i], actual[i]))
			return
		}
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
