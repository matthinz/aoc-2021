package d23

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseInputStarting(t *testing.T) {
	input := `
#############
#...........#
###C#A#B#D###
	#B#A#D#C#
	#########
		`
	g := parseInput(strings.NewReader(input))

	if g.hallway.height != 1 {
		t.Fatalf("Hallway height should not be %d", g.hallway.height)
	}

	if g.hallway.width != 11 {
		t.Fatalf("Hallway width should not be %d", g.hallway.width)
	}

	if len(g.rooms) != 4 {
		t.Fatalf("Should have 4 rooms, but have %d", len(g.rooms))
	}

	expectedRoomLocations := map[int]int{
		0: 2,
		1: 4,
		2: 6,
		3: 8,
	}
	for i, r := range g.rooms {
		if r.x != expectedRoomLocations[i] {
			t.Fatalf("Room %d should be at %d, but was at %d", i, expectedRoomLocations[i], r.x)
		}
	}

	expectedThings := map[int]amphipodKind{
		11: CopperAmphipod,
		12: BronzeAmphipod,
		13: AmberAmphipod,
		14: AmberAmphipod,
		15: BronzeAmphipod,
		16: DesertAmphipod,
		17: DesertAmphipod,
		18: CopperAmphipod,
	}

	for expectedPos, expectedKind := range expectedThings {
		found := false
		for a, actualPos := range g.initialState.positions {
			if actualPos == expectedPos {
				found = true
				if a.kind != expectedKind {
					t.Errorf("Expected position %d to have %s, but had %s", actualPos, string(expectedKind), string(a.kind))
				}
				break
			}
		}
		if !found {
			t.Errorf("Position %d not found (should have %s)", expectedPos, string(expectedKind))
		}
	}

}

func TestParseInputLater(t *testing.T) {
	input := `
#############
#C.........D#
###.#A#B#.###
	#B#A#D#C#
	#########
		`
	g := parseInput(strings.NewReader(input))

	if g.hallway.height != 1 {
		t.Fatalf("Hallway height should not be %d", g.hallway.height)
	}

	if g.hallway.width != 11 {
		t.Fatalf("Hallway width should not be %d", g.hallway.width)
	}

	if len(g.rooms) != 4 {
		t.Fatalf("Should have 4 rooms, but have %d", len(g.rooms))
	}

	expectedRoomLocations := map[int]int{
		0: 2,
		1: 4,
		2: 6,
		3: 8,
	}
	for i, r := range g.rooms {
		if r.x != expectedRoomLocations[i] {
			t.Fatalf("Room %d should be at %d, but was at %d", i, expectedRoomLocations[i], r.x)
		}
	}

	expectedThings := map[int]amphipodKind{
		0:  CopperAmphipod,
		10: DesertAmphipod,
		12: BronzeAmphipod,
		13: AmberAmphipod,
		14: AmberAmphipod,
		15: BronzeAmphipod,
		16: DesertAmphipod,
		18: CopperAmphipod,
	}

	for expectedPos, expectedKind := range expectedThings {
		found := false
		for a, actualPos := range g.initialState.positions {
			if actualPos == expectedPos {
				found = true
				if a.kind != expectedKind {
					t.Errorf("Expected position %d to have %s, but had %s", actualPos, string(expectedKind), string(a.kind))
				}
				break
			}
		}
		if !found {
			t.Errorf("Position %d not found (should have %s)", expectedPos, string(expectedKind))
		}
	}

}

func TestPositionInHallway(t *testing.T) {

	input := `
#############
#...........#
###C#A#B#D###
	#B#A#D#C#
	#########
		`
	g := parseInput(strings.NewReader(input))

	for x := 0; x < 11; x++ {
		pos := positionInHallway(&g, x)
		if pos != x {
			t.Errorf("positionInHallway wrong for %d (expected %d, got %d)", x, x, pos)
		}
	}
}

func TestPositionInRoom(t *testing.T) {
	input := `
#############
#...........#
###C#A#B#D###
	#B#A#D#C#
	#########
		`
	g := parseInput(strings.NewReader(input))

	tests := map[[2]int]int{
		{0, 0}: 11 + 0,
		{0, 1}: 11 + 1,
		{1, 0}: 11 + 2,
		{1, 1}: 11 + 3,
		{2, 0}: 11 + 4,
		{2, 1}: 11 + 5,
		{3, 0}: 11 + 6,
		{3, 1}: 11 + 7,
	}

	for input, expected := range tests {
		pos := positionInRoom(&g, input[0], input[1])
		if pos != expected {
			t.Errorf("Room %d, y %d returned the wrong position (got %d, expected %d)", input[0], input[1], pos, expected)
		}
	}

}

func TestApplyMove(t *testing.T) {
	input := `
#############
#...........#
###C#A#B#D###
	#B#A#D#C#
	#########
	`
	g := parseInput(strings.NewReader(input))

	m := move{
		from: 13,
		to:   0,
	}

	nextState := applyMove(&g, &g.initialState, m)

	expected := strings.TrimSpace(`
#############
#A..........#
###C#.#B#D###
###B#A#D#C###
#############
`)

	actual := stringify(&g, nextState)

	if actual != expected {
		t.Logf("EXPECTED\n%s", expected)
		t.Logf("ACTUAL\n%s", actual)
		t.Error("Move did not apply correctly")
	}
}

func TestFindLegalMovesForAmphipodInFirstPosition(t *testing.T) {

	input := `
#############
#...........#
###C#A#B#D###
	#B#A#D#C#
	#########
	`
	g := parseInput(strings.NewReader(input))

	var target *amphipod
	const targetPos = 13

	for a, pos := range g.initialState.positions {
		if pos == targetPos {
			if a.kind != AmberAmphipod {
				t.Fatalf("Expected pos %d to have amber, had %s", targetPos, string(a.kind))
			}
			target = a
			break
		}
	}

	if target == nil {
		t.Fatal("Could not find target")
	}

	ch := make(chan move)
	go func() {
		defer close(ch)
		findLegalMovesForAmphipod(&g, &g.initialState, target, targetPos, ch)
	}()

	moves := make([]move, 0)
	for m := range ch {
		fmt.Println(m)
		moves = append(moves, m)
	}

	expected := 10

	if len(moves) != expected {
		t.Errorf("Should've found %d moves, but found %d", expected, len(moves))
	}
}

func TestFindLegalMovesForAmphipodForGuysTuckedDeepInRooms(t *testing.T) {

	input := `
#############
#DB.......CA#
###.#.#.#.###
  #B#A#D#C#
  #########
	`
	g := parseInput(strings.NewReader(input))

	var target *amphipod
	const targetPos = 12

	fmt.Println(stringify(&g, &g.initialState))

	for a, pos := range g.initialState.positions {
		if pos == targetPos {
			if a.kind != BronzeAmphipod {
				t.Fatalf("Expected pos %d to have bronze, had %s", targetPos, string(a.kind))
			}
			target = a
			break
		}
	}

	if target == nil {
		t.Fatal("Could not find target")
	}

	ch := make(chan move)
	go func() {
		defer close(ch)
		findLegalMovesForAmphipod(&g, &g.initialState, target, targetPos, ch)
	}()

	moves := make([]move, 0)
	for m := range ch {
		fmt.Println(m)
		moves = append(moves, m)
	}

	expected := 1

	if len(moves) != expected {
		t.Errorf("Should've found %d moves, but found %d", expected, len(moves))
	}
}
