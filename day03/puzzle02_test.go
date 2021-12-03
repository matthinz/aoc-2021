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
	actual := FindO2GeneratorRating(input)
	expected := int64(23)
	if actual != expected {
		t.Error(fmt.Sprintf("Got the wrong rating! Expected %d, got %d", expected, actual))
	}
}

func TestFindCo2ScrubberRating(t *testing.T) {
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
	actual := FindCo2ScrubberRating(input)
	expected := int64(10)
	if actual != expected {
		t.Error(fmt.Sprintf("Got the wrong rating! Expected %d, got %d", expected, actual))
	}
}

func TestFindMostAndLeastCommonBits(t *testing.T) {
	input := []string{
		"1001",
		"1010",
		"0010",
		"1111",
	}

	expectedMostCommon := []int{
		1,
		0,
		1,
		1,
	}

	expectedLeastCommon := []int{
		0,
		1,
		0,
		0,
	}

	for i := 0; i < 4; i++ {
		mostCommon, leastCommon := FindMostAndLeastCommonBits(input, i)

		if mostCommon != expectedMostCommon[i] {
			t.Log(fmt.Sprintf("mostCommon should be %d at position %d", expectedMostCommon[i], i))
			t.Fail()
		}
		if leastCommon != expectedLeastCommon[i] {
			t.Log(fmt.Sprintf("leastCommon should be %d at position %d", expectedLeastCommon[i], i))
			t.Fail()
		}
	}
}

func TestFilter(t *testing.T) {
	input := []string{
		"1001",
		"1010",
		"0010",
		"1111",
	}

	actual := FilterNumbersByBitAtPosition(input, 1, 0)
	expected := []string{
		"1001",
		"1010",
		"1111",
	}

	assertEqual(t, actual, expected)

	actual = FilterNumbersByBitAtPosition(input, 0, 0)
	expected = []string{
		"0010",
	}
	assertEqual(t, actual, expected)

	actual = FilterNumbersByBitAtPosition(input, 0, 3)
	expected = []string{
		"1010",
		"0010",
	}

}

func assertEqual(t *testing.T, actual []string, expected []string) {

	if len(actual) != len(expected) {
		t.Error(actual)
	}
	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Error(fmt.Sprintf("%d: expected '%s', got '%s'", i, expected[i], actual[i]))
		}
	}

}
