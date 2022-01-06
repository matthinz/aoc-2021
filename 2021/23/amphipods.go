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
	hallwayWidth int
	roomHeight   int
	rooms        []room
	initialState gameState
}

type gameState struct {
	parent     *gameState
	lastMove   *move
	totalCost  int
	totalMoves int
	positions  map[*amphipod]int
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

	solvedState, statesEvaluated := solve(&g, &g.initialState, nil, l, 0)

	l.Printf("Evaluated %d total states", statesEvaluated)

	if solvedState == nil {
		panic("No solution found")
	}

	var moves []move
	for s := solvedState; s != nil; s = s.parent {
		if s.lastMove != nil {
			moves = append(moves, *s.lastMove)
		}
	}

	for i := len(moves) - 1; i >= 0; i-- {
		fmt.Printf("%d -> %d (%d)\n", moves[i].from, moves[i].to, moves[i].cost)
	}

	return strconv.Itoa(solvedState.totalCost)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// Part 1 solution

func solve(g *game, state *gameState, bestSolution *gameState, l *log.Logger, statesEvaluated uint) (*gameState, uint) {
	moves := getLegalMoves(g, state)
	for move := range moves {
		statesEvaluated++

		if bestSolution != nil && state.totalCost+move.cost >= bestSolution.totalCost {
			continue
		}

		nextState := applyMove(g, state, move)
		if isSolved(g, nextState) {
			if bestSolution == nil || nextState.totalCost < bestSolution.totalCost {
				bestSolution = nextState
				l.Printf("*** New best solution: %d (%d moves)", bestSolution.totalCost, bestSolution.totalMoves)
			}
			continue
		}

		solution, substatesEvaluated := solve(g, nextState, bestSolution, l, 0)
		statesEvaluated += substatesEvaluated
		if solution != nil {
			if bestSolution == nil || solution.totalCost < bestSolution.totalCost {
				bestSolution = solution
			}
		}
	}

	return bestSolution, statesEvaluated
}

func applyMove(g *game, state *gameState, m move) *gameState {
	nextState := gameState{
		parent:     state,
		lastMove:   &m,
		totalCost:  state.totalCost + m.cost,
		totalMoves: state.totalMoves + 1,
		positions:  make(map[*amphipod]int),
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
			if a.kind == AmberAmphipod {
				findLegalMovesForAmphipod(g, state, a, pos, ch)
			}
		}
		for a, pos := range state.positions {
			if a.kind == BronzeAmphipod {
				findLegalMovesForAmphipod(g, state, a, pos, ch)
			}
		}
		for a, pos := range state.positions {
			if a.kind == CopperAmphipod {
				findLegalMovesForAmphipod(g, state, a, pos, ch)
			}
		}
		for a, pos := range state.positions {
			if a.kind == DesertAmphipod {
				findLegalMovesForAmphipod(g, state, a, pos, ch)
			}
		}

	}()

	return ch
}

func findLegalMovesForAmphipod(g *game, state *gameState, a *amphipod, pos int, ch chan move) {

	// generate a slice mapping position indices to amphipods contained
	// so we can easily probe by position
	positions := make(
		[]*amphipod,
		g.hallwayWidth+(len(g.rooms)*g.roomHeight),
	)
	for a, pos := range state.positions {
		positions[pos] = a
	}

	startingPos := pos

	// First, try to move the amphipod from a room out into the hallway

	movedIntoHallway, hallwayPos, moveToHallwayCost := tryMoveAmphipodFromRoomToHallway(g, a, pos, positions)

	// Update our temporary position map
	if movedIntoHallway {
		positions[pos] = nil
		positions[hallwayPos] = a
		pos = hallwayPos
	}

	// Then try to move it from hallway to its destination room

	movedIntoDestinationRoom, destRoomPos, moveToDestRoomCost := tryMoveAmphipodFromHallwayToDestination(g, a, pos, positions)

	if movedIntoDestinationRoom {
		// We made it to the best place in the destination room, and this is the only move that matters.
		ch <- move{
			from: startingPos,
			to:   destRoomPos,
			cost: moveToHallwayCost + moveToDestRoomCost,
		}
		return
	}

	if movedIntoHallway {
		// We need to move this amphipod to a valid place in the hallway, otherwise it won't stick
		accessibleHallwayPositions := tryMoveAmphipodToValidPositionInHallway(g, a, hallwayPos, positions)
		for _, newHallwayPos := range accessibleHallwayPositions {
			ch <- move{
				from: startingPos,
				to:   newHallwayPos,
				cost: moveToHallwayCost + costToMove(a, int(math.Abs(float64(hallwayPos-newHallwayPos)))),
			}
		}
		return
	}

}

// attempts to move <a> from <hallwayPos> to another position in the hallway
// returns a slice of positions in the hallway that <a> can move to
func tryMoveAmphipodToValidPositionInHallway(g *game, a *amphipod, pos int, positions []*amphipod) []int {

	hallwayX, isInHallway := positionToHallwayX(g, pos)

	var accessiblePositions []int

	if !isInHallway {
		return accessiblePositions
	}

	canMoveLeft := true
	canMoveRight := true

	for deltaX := 1; deltaX < g.hallwayWidth; deltaX++ {
		if !(canMoveLeft || canMoveRight) {
			break
		}

		// try moving left
		if canMoveLeft {
			probeX := hallwayX - deltaX

			if probeX >= 0 {

				isAtRoomEntrance := false
				for i := range g.rooms {
					if g.rooms[i].x == probeX {
						isAtRoomEntrance = true
						break
					}
				}
				// amphipods are not allowed to stop right outside a room
				if !isAtRoomEntrance {

					probePos := positionInHallway(g, probeX)
					atProbePos := positions[probePos]
					if atProbePos == nil {
						accessiblePositions = append(accessiblePositions, probePos)
					} else {
						canMoveLeft = false
					}
				}
			}
		}
		// try moving right
		if canMoveRight {
			probeX := hallwayX + deltaX
			if probeX < g.hallwayWidth {

				isAtRoomEntrance := false
				for i := range g.rooms {
					if g.rooms[i].x == probeX {
						isAtRoomEntrance = true
						break
					}
				}

				if !isAtRoomEntrance {
					probePos := positionInHallway(g, probeX)
					atProbePos := positions[probePos]
					if atProbePos == nil {
						accessiblePositions = append(accessiblePositions, probePos)
					} else {
						canMoveRight = false
					}
				}
			}
		}
	}

	return accessiblePositions
}

// attempts to move an amphipod out of a room and into a hallway.
// returns flag indicating success, new hallway position, and cost of move into hallway
func tryMoveAmphipodFromRoomToHallway(g *game, a *amphipod, pos int, positions []*amphipod) (bool, int, int) {

	currentRoomIndex, currentRoomY, currentlyInAnyRoom := positionToRoomAndY(g, pos)

	if !currentlyInAnyRoom {
		// we can't move into a hallway if we're not in a room
		return false, 0, 0
	}

	destRoomIndex := getDestinationRoomIndex(a)

	if currentRoomIndex == destRoomIndex {
		// <a> is in their a destination room.
		// It will not need to move into the hallway unless it is blocking someone
		// else from getting out of the room

		isBlockingSomebody := false

		for y := currentRoomY + 1; y < g.roomHeight; y++ {
			probePos := positionInRoom(g, currentRoomIndex, y)
			atProbePos := positions[probePos]

			if atProbePos != nil && atProbePos.kind != a.kind {
				// <a> is blocking another amphipod, and must move out to the hallway
				isBlockingSomebody = true
				break
			}
		}

		if !isBlockingSomebody {
			// <a> is in its destination and is not blocking anybody, so it does
			// not need to move into the hallway
			return false, 0, 0
		}
	}

	// At this point, <a> needs to move into the hallway--either to advance to
	// its destination room or to clear the way for someone else to get out of
	// its current room.

	theWayOutIsBlocked := false
	for y := currentRoomY - 1; y >= 0; y-- {
		yPos := positionInRoom(g, currentRoomIndex, y)
		atY := positions[yPos]
		if atY != nil && atY != a {
			theWayOutIsBlocked = true
			break
		}
	}

	if theWayOutIsBlocked {
		// <a> cannot exit the room because someone else is in the way
		return false, 0, 0
	}

	// <a> can move out of its current room and into the hallway!
	// we move it out to the space immediately outside its room. this is not a
	// legal position in the long run, but that is a problem for some other
	// function.

	exitRoomCost := costToMove(a, currentRoomY+1)
	hallwayPos := positionInHallway(g, g.rooms[currentRoomIndex].x)

	return true, hallwayPos, exitRoomCost
}

// attempts to move an amphipod in the hallway into its destination room
// returns whether the move succeeded, the resulting position, and the cost
func tryMoveAmphipodFromHallwayToDestination(g *game, a *amphipod, pos int, positions []*amphipod) (bool, int, int) {

	x, isInHallway := positionToHallwayX(g, pos)

	if !isInHallway {
		return false, 0, 0
	}

	destRoomIndex := getDestinationRoomIndex(a)
	destRoomEntranceX := g.rooms[destRoomIndex].x

	step := 1
	if x > destRoomEntranceX {
		// need to move to the left
		step = -1
	}

	canReachDestinationRoom := true
	for probeX := x + step; probeX != destRoomEntranceX; probeX += step {
		if x < 0 || x >= g.hallwayWidth {
			break
		}
		probePos := positionInHallway(g, probeX)
		atPos := positions[probePos]

		if atPos != nil {
			// Someone is in the way
			canReachDestinationRoom = false
			break
		}
	}

	if !canReachDestinationRoom {
		return false, 0, 0
	}

	// Ok, we can reach it, but can we enter it?

	otherKindOfAmphipodInRoom := false
	shallowestBlocker := g.roomHeight

	// Move through the room and look for amphipods of any other kind as well
	// as amphipods of the same kind that may be blocking the way
	for probeY := g.roomHeight - 1; probeY >= 0; probeY-- {
		probePos := positionInRoom(g, destRoomIndex, probeY)
		atProbePos := positions[probePos]
		if atProbePos == nil {
			continue
		}

		if atProbePos.kind == a.kind {
			shallowestBlocker = probeY
		} else {
			otherKindOfAmphipodInRoom = true
			break
		}
	}

	if otherKindOfAmphipodInRoom {
		return false, 0, 0
	}

	if shallowestBlocker == 0 {
		// can't enter the room because it's blocked by an amphipod of the same kind
		return false, 0, 0
	}

	destPos := positionInRoom(g, destRoomIndex, shallowestBlocker-1)

	costToGetToRoomEntrance := costToMove(a, int(math.Abs(float64(destRoomEntranceX-x))))
	costToEnterRoom := costToMove(a, shallowestBlocker)

	return true, destPos, costToGetToRoomEntrance + costToEnterRoom
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
		g.hallwayWidth+(len(g.rooms)*g.roomHeight),
	)

	for a, pos := range state.positions {
		positions[pos] = a
	}

	for i := 0; i < g.hallwayWidth+2; i++ {
		b.WriteRune(wall)
	}
	b.WriteRune('\n')

	b.WriteRune(wall)
	for x := 0; x < g.hallwayWidth; x++ {
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

	for y := 0; y < g.roomHeight; y++ {

		b.WriteRune(wall)

		for x := 0; x < g.hallwayWidth; x++ {
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

	for x := 0; x < g.hallwayWidth+2; x++ {
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
	return g.hallwayWidth + (roomIndex * g.roomHeight) + y
}

func positionInHallway(g *game, x int) int {
	return x
}

func positionToRoomAndY(g *game, pos int) (int, int, bool) {
	if pos < g.hallwayWidth {
		return 0, 0, false
	}

	pos -= g.hallwayWidth

	y := pos % g.roomHeight
	roomIndex := (pos - y) / g.roomHeight

	return roomIndex, y, true
}

func positionToHallwayX(g *game, pos int) (int, bool) {
	if pos < g.hallwayWidth {
		return pos, true
	}
	return 0, false
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) game {

	s := bufio.NewScanner(r)

	g := game{
		hallwayWidth: 0,
		roomHeight:   2,
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
			for _, r := range line {
				if r == '.' {
					// empty slot in hallway
					g.hallwayWidth++
				}

				if r == rune(AmberAmphipod) || r == rune(BronzeAmphipod) || r == rune(CopperAmphipod) || r == rune(DesertAmphipod) {
					// an amphipod in the hallway
					g.hallwayWidth++
					a := amphipod{
						kind: amphipodKind(r),
					}
					// TODO: 2d hallway
					g.initialState.positions[&a] = g.hallwayWidth - 1
				}

			}
			continue
		}

		if g.hallwayWidth == 0 {
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
