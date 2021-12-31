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

func TestRollQuantumDie(t *testing.T) {
	expected := make(map[int]uint)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				expected[i+j+k] += 1
			}
		}
	}

	actual := rollQuantumDie(3, 3)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d moves, but got %d", len(expected), len(actual))
	}

	for move, universes := range expected {
		if actual[move] != universes {
			t.Errorf("Expected move %d to have %d universes, but had %d", move, universes, actual[move])
		}
	}

}
