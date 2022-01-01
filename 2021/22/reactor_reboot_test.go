package d22

import (
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	input := `
	on x=10..12,y=10..12,z=10..12
	on x=11..13,y=11..13,z=11..13
	off x=9..11,y=9..11,z=9..11
	on x=10..10,y=10..10,z=10..10
	`

	steps := parseInput(strings.NewReader(input))

	if len(steps) != 4 {
		t.Errorf("Expected %d steps, but got %d", 4, len(steps))
	}

	s := steps[0]
	if !s.turnOn {
		t.Error("First step should turn on")
	}

	if s.cuboid.size.x != 3 {
		t.Errorf("First step cuboid should have x dimension 3, but was %d", s.cuboid.size.x)
	}
	if s.cuboid.size.y != 3 {
		t.Errorf("First step cuboid should have y dimension 3, but was %d", s.cuboid.size.y)
	}
	if s.cuboid.size.z != 3 {
		t.Errorf("First step cuboid should have z dimension 3, but was %d", s.cuboid.size.z)
	}
}

func TestInitializationWithBruteForce(t *testing.T) {
	input := `
	on x=10..12,y=10..12,z=10..12
	on x=11..13,y=11..13,z=11..13
	off x=9..11,y=9..11,z=9..11
	on x=10..10,y=10..10,z=10..10
`
	steps := parseInput(strings.NewReader(input))

	reactor := applyStepsUsingBruteForce(steps)

	ct := countCubesOnUsingBruteForce(reactor)

	expected := uint(39)

	if ct != expected {
		t.Errorf("Expected %d, but got %d", expected, ct)
	}

}
