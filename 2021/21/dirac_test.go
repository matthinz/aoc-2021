package d21

import (
	"log"
	"strings"
	"testing"
)

func TestPuzzle1(t *testing.T) {
	input := strings.TrimSpace(`
Player 1 starting position: 4
Player 2 starting position: 8
`)
	actual := Puzzle1(strings.NewReader(input), log.Default())
	expected := "739785"
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestPuzzle2(t *testing.T) {
	input := strings.TrimSpace(`
Player 1 starting position: 4
Player 2 starting position: 8
`)
	actual := Puzzle2(strings.NewReader(input), log.Default())
	expected := "444356092776315"
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
