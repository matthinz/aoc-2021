package day04

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type square struct {
	value  int
	marked bool
}

type board struct {
	index   int
	squares [][]square
}

type game struct {
	numbers []int
	boards  []board
}

type solvedBoard struct {
	board
	lastNumberDrawn int
	finalScore      int
}

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(4, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	game := NewGame(r)
	solvedBoards := game.Run()

	board := solvedBoards[0]

	return strconv.Itoa(board.score())
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	game := NewGame(r)
	solvedBoards := game.Run()

	board := solvedBoards[len(solvedBoards)-1]

	return strconv.Itoa(board.score())
}

func main() {
	game := NewGame(os.Stdin)
	solvedBoards := game.Run()

	for _, sb := range solvedBoards {
		fmt.Printf("Board: %d/%d\nScore:%d\n", sb.index, len(solvedBoards), sb.finalScore)
	}
}

func (g *game) Run() []solvedBoard {

	var result []solvedBoard

	unsolvedBoards := g.boards

	for _, number := range g.numbers {
		for i := 0; i < len(unsolvedBoards); i++ {
			b := &unsolvedBoards[i]
			b.markNumber(number)

			if b.isSolved() {

				finalScore := b.sumOfUnmarkedSquares() * number

				// fmt.Printf("%d: %d * %d = %d\n", b.index, b.sumOfUnmarkedSquares(), number, finalScore)

				result = append(result, solvedBoard{
					board:           *b,
					lastNumberDrawn: number,
					finalScore:      finalScore,
				})

				unsolvedBoards = removeBoardAt(unsolvedBoards, i)
				i--
			}
		}
	}

	return result
}

func (b *board) String() string {
	result := strings.Builder{}

	for _, row := range b.squares {
		for _, square := range row {
			if square.marked {
				result.WriteString("(")
			} else {
				result.WriteString(" ")
			}
			if square.value < 10 {
				result.WriteString(" ")
			}
			result.WriteString(strconv.FormatInt(int64(square.value), 10))
			if square.marked {
				result.WriteString(")")
			} else {
				result.WriteString(" ")
			}
		}
		result.WriteString("\n")
	}

	return result.String()

}

func (b *board) markNumber(number int) {
	for y := range b.squares {
		for x := range b.squares[y] {
			s := &b.squares[y][x]
			if s.value == number {
				s.marked = true
			}
		}
	}
}

func (b *board) isSolved() bool {
	if len(b.squares) == 0 {
		return false
	}

	height := len(b.squares)
	width := len(b.squares[0])

	rowStates := make([]bool, height)
	colStates := make([]bool, width)

	// initialize state maps to true

	for y := 0; y < height; y++ {
		rowStates[y] = true
	}

	for x := 0; x < width; x++ {
		colStates[x] = true
	}

	for y, row := range b.squares {
		for x, square := range row {
			rowStates[y] = rowStates[y] && square.marked
			colStates[x] = colStates[x] && square.marked
		}
	}

	var anyRows bool
	var anyCols bool

	for _, ok := range rowStates {
		if ok {
			anyRows = true
		}
	}

	for _, ok := range colStates {
		if ok {
			anyCols = true
		}
	}

	return anyRows || anyCols
}

func (b *solvedBoard) score() int {
	return b.sumOfUnmarkedSquares() * b.lastNumberDrawn
}

func (b *board) sumOfUnmarkedSquares() int {
	var result int

	for _, row := range b.squares {
		for _, square := range row {
			if !square.marked {
				result += square.value
			}

		}
	}

	return result
}

func NewGame(r io.Reader) game {

	scanner := bufio.NewScanner(r)

	var g *game
	var b *board

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			b = nil
			continue
		}

		if g == nil {
			// first line of input = numbers
			g = &game{
				numbers: parseNumbers(line, ","),
			}
			continue
		}

		// subsequent lines of numbers = board
		if b == nil {
			g.boards = append(g.boards, board{
				index: len(g.boards) + 1,
			})
			b = &g.boards[len(g.boards)-1]
		}

		b.squares = append(b.squares, parseRow(line))
	}

	return *g
}

func parseNumbers(input string, sep string) []int {
	var result []int

	for _, token := range strings.Split(input, sep) {
		value, err := strconv.ParseInt(token, 10, 32)
		if err == nil {
			result = append(result, int(value))
		}
	}

	return result
}

func parseRow(input string) []square {
	var result []square
	for _, value := range parseNumbers(input, " ") {
		result = append(result, square{
			value: value,
		})
	}
	return result
}

func removeBoardAt(boards []board, index int) []board {
	var result []board
	for i := range boards {
		if i != index {
			result = append(result, boards[i])
		}
	}
	return result
}
