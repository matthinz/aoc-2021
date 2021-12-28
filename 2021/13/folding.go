package d13

import (
	_ "embed"
	"io"
	"log"
	"strconv"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func Puzzle1(r io.Reader, l *log.Logger) string {
	sheet := parseInput(r)
	for len(sheet.instructions) > 0 {
		i := sheet.instructions[0]
		nextSheet := sheet.fold()

		axis := "x"
		value := i.x
		if i.y > 0 {
			axis = "y"
			value = i.y
		}

		l.Printf("fold along %s=%d leaves %d points\n", axis, value, len(nextSheet.dots))

		return strconv.Itoa(len(nextSheet.dots))
	}

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	sheet := parseInput(r)
	for len(sheet.instructions) > 0 {
		i := sheet.instructions[0]
		nextSheet := sheet.fold()

		axis := "x"
		value := i.x
		if i.y > 0 {
			axis = "y"
			value = i.y
		}

		l.Printf("fold along %s=%d leaves %d points\n", axis, value, len(nextSheet.dots))
		sheet = nextSheet
	}

	return sheet.String()
}

func New() aoc.Day {
	return aoc.NewDay(13, defaultInput, Puzzle1, Puzzle2)
}
