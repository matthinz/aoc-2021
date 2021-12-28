package main

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestStep(t *testing.T) {
	input := strings.NewReader(strings.TrimSpace(`
	NNCB

	CH -> B
	HH -> N
	CB -> H
	NH -> C
	HB -> C
	HC -> B
	HN -> C
	NN -> C
	BH -> H
	NC -> B
	NB -> B
	BN -> B
	BB -> N
	BC -> B
	CC -> N
	CN -> C
	`))

	game := parseInput(input)

	result := run(game, 1)

	expected := "B=2,C=2,H=1,N=2,"
	actual := niceCounts(result)
	if actual != expected {
		t.Fatalf("Step 1 failed. Expected %v got %v", expected, actual)
	}
}

func Test4Steps(t *testing.T) {
	input := strings.NewReader(strings.TrimSpace(`
	NNCB

	CH -> B
	HH -> N
	CB -> H
	NH -> C
	HB -> C
	HC -> B
	HN -> C
	NN -> C
	BH -> H
	NC -> B
	NB -> B
	BN -> B
	BB -> N
	BC -> B
	CC -> N
	CN -> C
	`))

	game := parseInput(input)

	result := run(game, 4)

	expected := "B=23,C=10,H=5,N=11,"
	actual := niceCounts(result)
	if actual != expected {
		t.Fatalf("Step 4 failed. Expected %v got %v", expected, actual)
	}
}

func niceCounts(counts map[rune]int) string {

	sortedChars := make([]rune, 0, len(counts))
	for r := range counts {
		sortedChars = append(sortedChars, r)
	}
	sort.Slice(sortedChars, func(i, j int) bool {
		return sortedChars[i] < sortedChars[j]
	})

	result := strings.Builder{}
	for _, r := range sortedChars {
		result.WriteRune(r)
		result.WriteString("=")
		result.WriteString(strconv.Itoa(counts[r]))
		result.WriteString(",")
	}
	return result.String()
}
