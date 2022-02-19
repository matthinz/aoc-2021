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

	const Digits = 14

	l.Printf("parsing...")

	reg := parseInput(r)

	l.Printf("parsed")

	inputs, err := reg.z.FindInputs(0, PreferInputsThatMakeLargerNumber, l)

	if err != nil {
		panic(err)
	}

	digits := make([]int, Digits)

	for i := 0; i < Digits; i++ {
		digit, isSet := inputs[i]
		if isSet {
			digits[i] = digit
		} else {
			digits[i] = 9
		}
	}

	l.Println(digits)

	result := reg.z.Evaluate(digits)
	l.Println(result)

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}
