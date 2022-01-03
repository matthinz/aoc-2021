package d22

import (
	"log"
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

	cuboids := parseInput(strings.NewReader(input))

	if len(cuboids) != 4 {
		t.Errorf("Expected %d cuboids, but got %d", 4, len(cuboids))
	}

	c := cuboids[0]
	if !c.on {
		t.Error("First step should turn on")
	}

	if c.size.x != 3 {
		t.Errorf("First step cuboid should have x dimension 3, but was %d", c.size.x)
	}
	if c.size.y != 3 {
		t.Errorf("First step cuboid should have y dimension 3, but was %d", c.size.y)
	}
	if c.size.z != 3 {
		t.Errorf("First step cuboid should have z dimension 3, but was %d", c.size.z)
	}
}

func TestBuildIntervals(t *testing.T) {
	input := `
	on x=10..12,y=10..12,z=10..12
	on x=11..13,y=11..13,z=11..13
`
	cuboids := parseInput(strings.NewReader(input))

	xIntervals := buildIntervals(
		cuboids,
		func(c cuboid) int {
			return c.position.x
		},
		func(c cuboid) int {
			return c.size.x
		},
	)

	expected := []interval{
		{
			start:         10,
			end:           11,
			cuboidIndices: []int{0},
		},
		{
			start:         11,
			end:           13,
			cuboidIndices: []int{0, 1},
		},
		{
			start:         13,
			end:           14,
			cuboidIndices: []int{1},
		},
	}

	if len(xIntervals) != len(expected) {
		t.Fatalf("Wrong # of intervals. Expected %v (%d), got %v (%d)", expected, len(expected), xIntervals, len(xIntervals))
	}

	for i, actual := range xIntervals {

		if actual.start != expected[i].start {
			t.Log(xIntervals)
			t.Fatalf("#%d has bad start. Expected %d, got %d", i, expected[i].start, actual.start)
		}

		if actual.end != expected[i].end {
			t.Log(xIntervals)
			t.Fatalf("#%d has bad end. Expected %d, got %d", i, expected[i].end, actual.end)
		}

		if len(actual.cuboidIndices) != len(expected[i].cuboidIndices) {
			t.Log(xIntervals)
			t.Fatalf("#%d has wrong # of cuboidIndices. Expected %v (%d), got %v, (%d)", i, expected[i].cuboidIndices, len(expected[i].cuboidIndices), actual.cuboidIndices, len(actual.cuboidIndices))
		} else {
			ok := true
			for j := range actual.cuboidIndices {
				if actual.cuboidIndices[j] != expected[i].cuboidIndices[j] {
					ok = false
				}
			}
			if !ok {
				t.Log(xIntervals)
				t.Fatalf("#%d has wrong cuboid indices. Expected %v (%d), got %v, (%d)", i, expected[i].cuboidIndices, len(expected[i].cuboidIndices), actual.cuboidIndices, len(actual.cuboidIndices))
			}
		}
	}
}

func TestInitializationWithoutBruteForce(t *testing.T) {
	input := `
	on x=10..12,y=10..12,z=10..12
	on x=11..13,y=11..13,z=11..13
	off x=9..11,y=9..11,z=9..11
	on x=10..10,y=10..10,z=10..10
`
	cuboids := parseInput(strings.NewReader(input))

	t.Logf("Parsed %d cuboids from input", len(cuboids))

	normalized := initializeReactor(cuboids, log.Default())

	t.Logf("Normalized into %d cuboids", len(normalized))

	ct := countCubesOn(normalized)

	expected := uint(39)

	if ct != expected {
		t.Errorf("Expected %d, but got %d", expected, ct)
	}
}
