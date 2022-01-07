package d23

import (
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

	if g.hallwayWidth != 11 {
		t.Fatalf("Hallway width should not be %d", g.hallwayWidth)
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
		for actualPos, a := range g.initialState.positions {
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

	if g.hallwayWidth != 11 {
		t.Fatalf("Hallway width should not be %d", g.hallwayWidth)
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
		for actualPos, a := range g.initialState.positions {
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

	for pos, a := range g.initialState.positions {
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
		moves = append(moves, m)
	}

	expected := 7

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

	// fmt.Println(stringify(&g, &g.initialState))

	for pos, a := range g.initialState.positions {
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
		moves = append(moves, m)
	}

	expected := 3

	if len(moves) != expected {
		t.Errorf("Should've found %d moves, but found %d", expected, len(moves))
	}
}

func TestRunThroughExample(t *testing.T) {

	type test struct {
		input string
		kind  amphipodKind
		from  int
		to    int
		cost  int
	}

	tests := []test{
		{
			input: `
				#############
				#...........#
				###B#C#B#D###
					#A#D#C#A#
					#########
			`,
			kind: BronzeAmphipod,
			from: 15,
			to:   3,
			cost: 40,
		},
		{
			input: `
				#############
				#...B.......#
				###B#C#.#D###
					#A#D#C#A#
					#########
			`,
			kind: CopperAmphipod,
			from: 13,
			to:   15,
			cost: 400,
		},
		{
			input: `
				#############
				#...B.......#
				###B#.#C#D###
					#A#D#C#A#
					#########
			`,
			kind: DesertAmphipod,
			from: 14,
			to:   5,
			cost: 3000,
		},
		{
			input: `
				#############
				#...B.D.....#
				###B#.#C#D###
					#A#.#C#A#
					#########
			`,
			kind: BronzeAmphipod,
			from: 3,
			to:   14,
			cost: 30,
		},
		{
			input: `
				#############
				#.....D.....#
				###B#.#C#D###
					#A#B#C#A#
					#########
			`,
			kind: BronzeAmphipod,
			from: 11,
			to:   13,
			cost: 40,
		},
		{
			input: `
				#############
				#.....D.....#
				###.#B#C#D###
					#A#B#C#A#
					#########
			`,
			kind: DesertAmphipod,
			from: 17,
			to:   7,
			cost: 2000,
		},
		{
			input: `
				#############
				#.....D.D...#
				###.#B#C#.###
					#A#B#C#A#
					#########
			`,
			kind: AmberAmphipod,
			from: 18,
			to:   9,
			cost: 3,
		},
		{
			input: `
				#############
				#.....D.D.A.#
				###.#B#C#.###
					#A#B#C#.#
					#########
			`,
			kind: DesertAmphipod,
			from: 7,
			to:   18,
			cost: 3000,
		},
		{
			input: `
				#############
				#.....D...A.#
				###.#B#C#.###
					#A#B#C#D#
					#########
			`,
			kind: DesertAmphipod,
			from: 5,
			to:   17,
			cost: 4000,
		},
		{
			input: `
				#############
				#.........A.#
				###.#B#C#D###
					#A#B#C#D#
					#########
			`,
			kind: AmberAmphipod,
			from: 9,
			to:   11,
			cost: 8,
		},
	}

	for testIndex, test := range tests {

		g := parseInput(strings.NewReader(test.input))

		var a *amphipod
		for pos, check := range g.initialState.positions {
			if pos == test.from {
				a = check
				break
			}
		}

		if a == nil {
			t.Errorf("Test %d: No amphipod found at position %d (expected %s)", testIndex, test.from, string(test.kind))
		}

		ch := make(chan move)
		go func() {
			defer close(ch)
			findLegalMovesForAmphipod(&g, &g.initialState, a, test.from, ch)
		}()

		foundMove := false

		for m := range ch {
			if m.from == test.from && m.to == test.to {
				if m.cost != test.cost {
					t.Errorf("Test %d: Found move from %d to %d, but cost was wrong (expected %d, got %d)", testIndex, test.from, test.to, test.cost, m.cost)
				}
				foundMove = true
				break
			}
		}

		if !foundMove {
			t.Errorf("Test %d: Did not find move from %d to %d", testIndex, test.from, test.to)
		}

	}

}

func TestIsSolvedSuccess(t *testing.T) {
	input := `
	#############
	#...........#
	###A#B#C#D###
		#A#B#C#D#
		#########
	`

	g := parseInput(strings.NewReader(input))

	if !isSolved(&g, &g.initialState) {
		t.Errorf("Should detect solved state")
	}

}

func TestUnfoldDiagram(t *testing.T) {

	input := `
#############
#...........#
###B#C#B#D###
  #A#D#C#A#
  #########
	`

	expected := strings.TrimSpace(`
#############
#...........#
###B#C#B#D###
  #D#C#B#A#
  #D#B#A#C#
  #A#D#C#A#
  #########
	`)

	unfolded := unfoldDiagram(strings.NewReader(input))

	if unfolded != expected {
		t.Log(unfolded)
		t.Fatalf("Unfold did not work!")
	}

}
