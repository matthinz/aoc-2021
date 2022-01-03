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

func TestInitializationWithBruteForce(t *testing.T) {
	input := `
	on x=10..12,y=10..12,z=10..12
	on x=11..13,y=11..13,z=11..13
	off x=9..11,y=9..11,z=9..11
	on x=10..10,y=10..10,z=10..10
`
	cuboids := parseInput(strings.NewReader(input))

	reactor := initializeReactorUsingBruteForce(cuboids)

	ct := countCubesOnUsingBruteForce(reactor)

	expected := uint(39)

	if ct != expected {
		t.Errorf("Expected %d, but got %d", expected, ct)
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
			end:           10,
			cuboidIndices: []int{0},
		},
		{
			start:         11,
			end:           12,
			cuboidIndices: []int{0, 1},
		},
		{
			start:         13,
			end:           13,
			cuboidIndices: []int{1},
		},
	}

	if len(xIntervals) != len(expected) {
		t.Fatalf("Wrong # of intervals. Expected %v (%d), got %v (%d)", expected, len(expected), xIntervals, len(xIntervals))
	}

	for i, actual := range xIntervals {

		if actual.start != expected[i].start {
			t.Errorf("#%d has bad start. Expected %d, got %d", i, expected[i].start, actual.start)
		}

		if actual.end != expected[i].end {
			t.Errorf("#%d has bad end. Expected %d, got %d", i, expected[i].end, actual.end)
		}

		if len(actual.cuboidIndices) != len(expected[i].cuboidIndices) {
			t.Errorf("#%d has wrong # of cuboidIndices. Expected %v (%d), got %v, (%d)", i, expected[i].cuboidIndices, len(expected[i].cuboidIndices), actual.cuboidIndices, len(actual.cuboidIndices))
		} else {
			ok := true
			for j := range actual.cuboidIndices {
				if actual.cuboidIndices[j] != expected[i].cuboidIndices[j] {
					ok = false
				}
			}
			if !ok {
				t.Errorf("#%d has wrong steps. Expected %v (%d), got %v, (%d)", i, expected[i].cuboidIndices, len(expected[i].cuboidIndices), actual.cuboidIndices, len(actual.cuboidIndices))
			}
		}
	}
}
