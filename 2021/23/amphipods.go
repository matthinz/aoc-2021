package d23

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type amphipodKind rune

type amphipod struct {
	kind amphipodKind
}

type room struct {
	position  int
	occupants []*amphipod
}

type game struct {
	hallway []*amphipod
	rooms   []room
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
	game := parseInput(r)
	fmt.Println(game.String())
	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// game

func (g *game) String() string {
	b := strings.Builder{}

	for i := 0; i < len(g.hallway)+2; i++ {
		b.WriteRune('#')
	}
	b.WriteRune('\n')

	b.WriteRune('#')
	for i := 0; i < len(g.hallway); i++ {
		if g.hallway[i] == nil {
			b.WriteRune('.')
		} else {
			b.WriteRune(rune(g.hallway[i].kind))
		}
	}
	b.WriteRune('#')
	b.WriteRune('\n')

	var maxRoomDepth, minRoomX, maxRoomX int

	for i := range g.rooms {
		r := &g.rooms[i]
		if len(r.occupants) > maxRoomDepth {
			maxRoomDepth = len(r.occupants)
		}
		if minRoomX == 0 || r.position < minRoomX {
			minRoomX = r.position
		}
		if maxRoomX == 0 || r.position > maxRoomX {
			maxRoomX = r.position
		}
	}

	for depth := 0; depth < maxRoomDepth; depth++ {

		padChar := ' '
		if depth == 0 {
			padChar = '#'
		}

		b.WriteRune(padChar)

		for x := 0; x < len(g.hallway); x++ {
			var roomAtPosition, roomAtPrevPosition, roomAtNextPosition *room

			for i := range g.rooms {
				r := &g.rooms[i]
				if r.position == x {
					roomAtPosition = r
				} else if r.position == x+1 {
					roomAtNextPosition = r
				} else if r.position == x-1 {
					roomAtPrevPosition = r
				}
			}

			if roomAtPosition != nil {
				isEmpty := len(roomAtPosition.occupants) < depth+1 || roomAtPosition.occupants[depth] == nil
				if isEmpty {
					b.WriteRune('.')
				} else {
					occupant := roomAtPosition.occupants[depth]
					b.WriteRune(rune(occupant.kind))
				}
				continue
			}

			if roomAtNextPosition != nil || roomAtPrevPosition != nil {
				b.WriteRune('#')
			} else {
				b.WriteRune(padChar)
			}
		}

		b.WriteRune(padChar)

		b.WriteRune('\n')
	}

	b.WriteRune(' ')

	for i := 0; i < len(g.hallway); i++ {
		if i >= minRoomX-1 && i <= maxRoomX+1 {
			b.WriteRune('#')
		} else {
			b.WriteRune(' ')
		}
	}

	b.WriteRune(' ')

	return b.String()
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) game {
	g := game{}

	s := bufio.NewScanner(r)

	var hallwayWidth int

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		if strings.Index(line, "#.") == 0 {
			for _, r := range line {
				if r == '.' {
					hallwayWidth++
				}
			}

			g.hallway = make([]*amphipod, hallwayWidth)

			continue
		}

		if hallwayWidth == 0 {
			continue
		}

		roomIndex := 0
		for i, r := range line {
			if r == '#' || r == ' ' {
				continue
			}
			switch amphipodKind(r) {
			case AmberAmphipod, BronzeAmphipod, CopperAmphipod, DesertAmphipod:

				// put this amphipod in a room
				for len(g.rooms) < roomIndex+1 {
					g.rooms = append(g.rooms, room{
						position: i - 1,
					})
				}

				a := amphipod{
					kind: amphipodKind(r),
				}

				g.rooms[roomIndex].occupants = append(g.rooms[roomIndex].occupants, &a)
				roomIndex++

			default:
				// ignore

			}
		}
	}

	fmt.Println(g)

	return g
}
