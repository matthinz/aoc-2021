package d24

import (
	_ "embed"
	"io"
	"log"

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

	inputs, err := Solve(reg.z, 0, l)

	if err != nil {
		panic(err)
	}

	l.Print(inputs)

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}
