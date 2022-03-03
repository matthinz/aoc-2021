package d24

import (
	_ "embed"
	"io"
	"log"
	"strconv"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(24, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	l.Printf("parsing...")
	reg := parseInput(r)
	l.Printf("parse completed")

	inputs, err := SolveForLargest(reg.z, 0, l)

	if err != nil {
		panic(err)
	}

	result := ""
	for i := 0; i < 14; i++ {
		result += strconv.Itoa(inputs[i])
	}

	return result
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	l.Printf("parsing...")
	reg := parseInput(r)
	l.Printf("parse completed")

	inputs, err := SolveForSmallest(reg.z, 0, l)

	if err != nil {
		panic(err)
	}

	result := ""
	for i := 0; i < 14; i++ {
		result += strconv.Itoa(inputs[i])
	}

	return result
}
