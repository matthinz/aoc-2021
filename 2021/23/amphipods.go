package d23

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type amphipodKind rune

type amphipod struct {
	kind amphipodKind
}

type size struct {
	width, height int
}

type room struct {
	x int
}

type game struct {
	hallway      size
	rooms        []room
	roomSize     size
	initialState gameState
}

type gameState struct {
	parent    *gameState
	lastMove  *move
	totalCost int
	positions map[*amphipod]int
}

type move struct {
	from, to int
	cost     int
}

const (
	AmberAmphipod  amphipodKind = 'A'
	BronzeAmphipod              = 'B'
	CopperAmphipod              = 'C'
	DesertAmphipod              = 'D'
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(23, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	g := parseInput(r)

	solvedState := solve(&g, &g.initialState, nil)

	if solvedState == nil {
		panic("No solution found")
	}

	return strconv.Itoa(solvedState.totalCost)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// Part 1 solution

func solve(g *game, state *gameState, bestSolution *gameState) *gameState {
	moves := getLegalMoves(g, state)
	for move := range moves {
		nextState := applyMove(g, state, move)
		fmt.Println(stringify(g, nextState))
		if isSolved(g, nextState) {
			if bestSolution == nil || nextState.totalCost < bestSolution.totalCost {
				bestSolution = nextState
			}
			continue
		}

		solution := solve(g, nextState, bestSolution)
		if solution != nil {
			if bestSolution == nil || solution.totalCost < bestSolution.totalCost {
				bestSolution = solution
			}
		}
	}

	return bestSolution
}

func applyMove(g *game, state *gameState, m move) *gameState {
	nextState := gameState{
		parent:    state,
		lastMove:  &m,
		totalCost: state.totalCost + m.cost,
		positions: make(map[*amphipod]int),
	}

	var from *amphipod

	for a, pos := range state.positions {
		if pos == m.from {
			from = a
		} else {
			nextState.positions[a] = pos
		}
	}

	if from == nil {
		panic("can't move nothing?")
	}

	nextState.positions[from] = m.to

	return &nextState
}

func costToMove(a *amphipod, spaces int) int {
	switch a.kind {
	case AmberAmphipod:
		return spaces

	case BronzeAmphipod:
		return spaces * 10

	case CopperAmphipod:
		return spaces * 100

	case DesertAmphipod:
		return spaces * 1000
	default:
		panic(fmt.Sprintf("Unknown AmphipodKind: %s", string(a.kind)))
	}
}

func getDestinationRoomIndex(a *amphipod) int {

	switch a.kind {
	case AmberAmphipod:
		return 0
	case BronzeAmphipod:
		return 1
	case CopperAmphipod:
		return 2
	case DesertAmphipod:
		return 3
	default:
		panic(fmt.Sprintf("Unknown AmphipodKind: %s", string(a.kind)))
	}

}

func getLegalMoves(g *game, state *gameState) chan move {
	ch := make(chan move)

	go func() {
		defer close(ch)

		for a, pos := range state.positions {
			findLegalMovesForAmphipod(g, state, a, pos, ch)
		}
	}()

	return ch
}

func findLegalMovesForAmphipod(g *game, state *gameState, a *amphipod, pos int, ch chan move) {

	// generate a slice mapping position indices to amphipods contained
	// so we can easily probe by position
	positions := make(
		[]*amphipod,
		(g.hallway.width*g.hallway.height)+len(g.rooms)*(g.roomSize.width*g.roomSize.height),
	)
	for a, pos := range state.positions {
		positions[pos] = a
	}

	destRoomIndex := getDestinationRoomIndex(a)
	currentRoomIndex, currentRoomY, currentlyInAnyRoom := positionToRoomAndY(g, pos)

	startingPos := pos
	startedInARoom := currentlyInAnyRoom
	exitRoomCost := 0

	if currentlyInAnyRoom {

		if currentRoomIndex == destRoomIndex {
			// <a> is in their a destination room. if it is optimally placed
			// (can't move further into the room) and not obstructing anybody,
			// it does not need to move any more.
			deepestOpen := 0
			blockingSomebody := false

			for y := currentRoomY + 1; y < g.roomSize.height; y++ {
				atY := positions[positionInRoom(g, currentRoomIndex, y)]
				if atY == nil {
					// there is nobody at this position, so <a> could move deeper into
					// its room. this should not happen?
					deepestOpen = y
				} else if atY.kind != a.kind {
					// <a> is blocking another amphipod, and must move out to the hallway
					blockingSomebody = true
				}
			}

			if !blockingSomebody && deepestOpen == 0 {
				// <a> is optimally placed in its destination
				return
			}

			if !blockingSomebody {
				// <a> is in its destination, but could move deeper in
				ch <- move{
					from: pos,
					to:   positionInRoom(g, currentRoomIndex, deepestOpen),
					cost: costToMove(a, deepestOpen-currentRoomY),
				}
				return
			}

			// at this point, <a> is blocking someone, and so must move out into the
			// hallway. we fall through to handle that case
		}

		// <a> needs to move out of the room it is in.
		// make sure the way is clear
		wayIsBlocked := false
		for y := currentRoomY - 1; y >= 0; y++ {
			atY := positions[positionInRoom(g, currentRoomIndex, y)]
			if atY != nil {
				wayIsBlocked = true
				break
			}
		}

		if wayIsBlocked {
			// no moves can be made
			return
		}

		// we can move the amphipod *at least* to the hallway
		exitRoomCost += costToMove(a, currentRoomY+1)
		pos = positionInHallway(g, g.rooms[currentRoomIndex].x)
	}

	canEnterDestination := true
	shallowestBlocker := g.roomSize.height

	for y := g.roomSize.height - 1; y >= 0; y-- {
		probePos := positionInRoom(g, destRoomIndex, y)
		atPos := positions[probePos]
		if atPos == nil {
			continue
		}

		if atPos.kind == a.kind {
			if y < shallowestBlocker {
				shallowestBlocker = y
			}
		} else {
			canEnterDestination = false
			break
		}
	}

	if shallowestBlocker == 0 {
		// can't enter because it's plum full up
		canEnterDestination = false
	}

	if canEnterDestination {
		// try to move cleanly from our position in the hallway into the destination room

		hallwayX, currentlyInHallway := positionToHallwayX(g, pos)
		if !currentlyInHallway {
			panic("We're not in the hallway? But we should be?")
		}

		destPos := positionInRoom(g, destRoomIndex, shallowestBlocker-1)

		step := 1
		if hallwayX < g.rooms[destRoomIndex].x {
			step = -1
		}

		canReachDestinationRoom := true
		for x := hallwayX; x != g.rooms[destRoomIndex].x; x += step {
			probePos := positionInHallway(g, x)
			atPos := positions[probePos]
			if atPos != nil {
				canReachDestinationRoom = false
				break
			}
		}

		if !canReachDestinationRoom {
			return
		}

		// cost to move from current position in hallway to position immediately
		// above the room
		hallwayCost := costToMove(
			a,
			int(math.Abs(float64(hallwayX-g.rooms[destRoomIndex].x))),
		)

		enterRoomCost := costToMove(
			a,
			shallowestBlocker, // includes 1 step to move from hallway into room
		)

		ch <- move{
			from: startingPos,
			to:   destPos,
			cost: exitRoomCost + hallwayCost + enterRoomCost,
		}

		return
	}

	// When we can't enter the destination, if we started in the hallway, then
	// we can't actually do *anything* at all
	if !startedInARoom {
		return
	}

	// Ok, so we can't enter the destination, but we need to try to move
	// *somewhere* in the hallway.
	hallwayX, currentlyInHallway := positionToHallwayX(g, pos)
	if !currentlyInHallway {
		panic("We're not in the hallway? But we should be?")
	}

	canMoveLeft := true
	canMoveRight := true

	for deltaX := 1; deltaX < g.hallway.width; deltaX++ {
		if !(canMoveLeft || canMoveRight) {
			break
		}

		// try moving left
		if canMoveLeft {
			probeX := hallwayX - deltaX
			if probeX >= 0 {
				probePos := positionInHallway(g, probeX)
				atProbePos := positions[probePos]
				if atProbePos == nil {
					// nobody at this position, ok to move
					ch <- move{
						from: startingPos,
						to:   probePos,
						cost: exitRoomCost + costToMove(a, deltaX),
					}
				} else {
					canMoveLeft = false
				}
			}
		}
		// try moving right
		if canMoveRight {
			probeX := hallwayX + deltaX
			if probeX < g.hallway.width {
				probePos := positionInHallway(g, probeX)
				atProbePos := positions[probePos]
				if atProbePos == nil {
					// nobody at this position, ok to move
					ch <- move{
						from: startingPos,
						to:   probePos,
						cost: exitRoomCost + costToMove(a, deltaX),
					}
				} else {
					canMoveRight = false
				}
			}
		}
	}
}

func isSolved(g *game, state *gameState) bool {

	for a, pos := range state.positions {

		roomIndex, _, inRoomAtAll := positionToRoomAndY(g, pos)

		if !inRoomAtAll {
			return false
		}

		if roomIndex != getDestinationRoomIndex(a) {
			return false
		}
	}

	return true
}

////////////////////////////////////////////////////////////////////////////////
// game

func stringify(g *game, state *gameState) string {
	b := strings.Builder{}

	const wall = '#'
	const blank = ' '
	const open = '.'

	positions := make(
		[]*amphipod,
		(g.hallway.width*g.hallway.height)+len(g.rooms)*(g.roomSize.width*g.roomSize.height),
	)

	for a, pos := range state.positions {
		positions[pos] = a
	}

	for i := 0; i < g.hallway.width+2; i++ {
		b.WriteRune(wall)
	}
	b.WriteRune('\n')

	for y := 0; y < g.hallway.height; y++ {
		b.WriteRune(wall)
		for x := 0; x < g.hallway.width; x++ {
			pos := positionInHallway(g, x)
			atPos := positions[pos]
			if atPos == nil {
				b.WriteRune(open)
			} else {
				b.WriteRune(rune(atPos.kind))
			}
		}
		b.WriteRune(wall)
		b.WriteRune('\n')
	}

	for y := 0; y < g.roomSize.height; y++ {

		b.WriteRune(wall)

		for x := 0; x < g.hallway.width; x++ {
			var room *room
			for i := range g.rooms {
				if g.rooms[i].x == x {
					room = &g.rooms[i]
				}
			}

			if room != nil {

				a := amphipodInRoom(g, state, x, y)
				if a != nil {
					b.WriteRune(rune(a.kind))
				} else {
					b.WriteRune(open)
				}
				continue
			}

			b.WriteRune(wall)

		}

		b.WriteRune(wall)
		b.WriteRune('\n')
	}

	for x := 0; x < g.hallway.width+2; x++ {
		b.WriteRune(wall)
	}

	return b.String()
}

func amphipodInRoom(g *game, state *gameState, x, y int) *amphipod {

	roomIndex := -1
	for i := range g.rooms {
		if g.rooms[i].x == x {
			roomIndex = i
			break
		}
	}

	if roomIndex == -1 {
		return nil
	}

	pos := positionInRoom(g, roomIndex, y)

	for a, candidatePos := range state.positions {
		if candidatePos == pos {
			return a
		}
	}

	return nil
}

func positionInRoom(g *game, roomIndex int, y int) int {
	return (g.hallway.width * g.hallway.height) + (roomIndex * (g.roomSize.width * g.roomSize.height)) + y
}

func positionInHallway(g *game, x int) int {
	return x
}

func positionToRoomAndY(g *game, pos int) (int, int, bool) {
	hallwayArea := g.hallway.width * g.hallway.height
	if pos < hallwayArea {
		return 0, 0, false
	}

	pos -= hallwayArea

	roomArea := (g.roomSize.width * g.roomSize.height)
	y := pos % roomArea
	roomIndex := (pos - y) / roomArea

	return roomIndex, y, true
}

func positionToHallwayX(g *game, pos int) (int, bool) {
	hallwayArea := g.hallway.width * g.hallway.height
	if pos < hallwayArea {
		return pos, true
	}
	return 0, false
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) game {

	s := bufio.NewScanner(r)

	g := game{
		hallway:  size{},
		roomSize: size{width: 1, height: 2},
		initialState: gameState{
			positions: make(map[*amphipod]int),
		},
	}

	roomY := 0

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		isAllWalls, _ := regexp.MatchString("^#+$", line)
		if isAllWalls {
			continue
		}

		isRoomLine, _ := regexp.MatchString("^(\\s|#)*#((\\.|[ABCD])#)+", line)

		if !isRoomLine {
			g.hallway.height++
			for _, r := range line {
				if r == '.' {
					// empty slot in hallway
					g.hallway.width++
				}

				if r == rune(AmberAmphipod) || r == rune(BronzeAmphipod) || r == rune(CopperAmphipod) || r == rune(DesertAmphipod) {
					// an amphipod in the hallway
					g.hallway.width++
					a := amphipod{
						kind: amphipodKind(r),
					}
					// TODO: 2d hallway
					g.initialState.positions[&a] = g.hallway.width - 1
				}

			}
			continue
		}

		if g.hallway.width == 0 {
			continue
		}

		roomIndex := 0
		for i, r := range line {
			if r == '#' || r == ' ' {
				continue
			}

			if r == '.' {
				// this is an empty slot in the room
				if len(g.rooms) < roomIndex+1 {
					g.rooms = append(g.rooms, room{
						x: i - 1,
					})
				}
				roomIndex++
				continue
			}

			switch amphipodKind(r) {
			case AmberAmphipod, BronzeAmphipod, CopperAmphipod, DesertAmphipod:

				if len(g.rooms) < roomIndex+1 {
					g.rooms = append(g.rooms, room{
						x: i - 1,
					})
				}

				a := amphipod{
					kind: amphipodKind(r),
				}

				g.initialState.positions[&a] = positionInRoom(&g, roomIndex, roomY)

				roomIndex++

			default:
				panic(fmt.Sprintf("Unrecognized character in input: %s", string(r)))
			}
		}

		roomY++
	}

	return g
}
