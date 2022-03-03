package d25

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

type seaCucumber rune

const (
	noSeaCucumber    seaCucumber = '.'
	eastSeaCucumber              = '>'
	southSeaCucumber             = 'v'
)

func New() aoc.Day {
	return aoc.NewDay(25, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	board := parseInput(r)
	move := 0
	for {
		move++
		nextBoard, moved := tick(board)

		for _, row := range *nextBoard {
			for _, cell := range row {
				fmt.Print(string(cell))
			}
			fmt.Println()
		}
		fmt.Println(moved)

		if moved == 0 {
			return strconv.Itoa(move)
		}
		board = nextBoard
	}
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

func parseInput(r io.Reader) *[][]seaCucumber {
	result := make([][]seaCucumber, 0)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		row := make([]seaCucumber, len(line))
		for x, r := range line {
			switch r {
			case '.':
				row[x] = noSeaCucumber
			case 'v':
				row[x] = southSeaCucumber
			case '>':
				row[x] = eastSeaCucumber
			default:
				panic(fmt.Sprintf("invalid char: '%s'", string(r)))
			}
		}
		result = append(result, row)
	}

	return &result
}

func tick(board *[][]seaCucumber) (*[][]seaCucumber, int) {
	nextBoard := make([][]seaCucumber, len(*board))

	// 0. Prepare next board
	for y := range *board {
		row := (*board)[y]
		width := len(row)
		nextRow := make([]seaCucumber, width)
		if y == 0 {
			for x := 0; x < width; x++ {
				nextRow[x] = noSeaCucumber
			}
		} else {
			copy(nextRow, nextBoard[0])
		}
		nextBoard[y] = nextRow
	}

	moved := 0

	// 1. process east-facing
	for y := range *board {
		row := (*board)[y]
		width := len(row)
		for x := range row {
			cell := row[x]

			if cell != eastSeaCucumber {
				continue
			}

			destX := x + 1
			if destX >= width {
				destX = 0
			}

			atDest := row[destX]
			if atDest == noSeaCucumber {
				nextBoard[y][destX] = eastSeaCucumber
				moved++
			} else {
				nextBoard[y][x] = eastSeaCucumber
			}
		}
	}

	// 2. Process south-facing
	for y := range *board {
		row := (*board)[y]
		for x := range row {
			cell := row[x]
			if cell != southSeaCucumber {
				continue
			}

			destY := y + 1
			if destY >= len(*board) {
				destY = 0
			}

			// We have to make sure there's no south-facing cucumber at our destination
			// and no east-facing cucumber at the destination on the next

			atDestOnCurrent := (*board)[destY][x]
			atDestOnNext := nextBoard[destY][x]

			if atDestOnCurrent == southSeaCucumber || atDestOnNext != noSeaCucumber {
				// we can't move
				nextBoard[y][x] = southSeaCucumber
				continue
			}

			nextBoard[destY][x] = southSeaCucumber
			moved++
		}
	}

	return &nextBoard, moved
}
