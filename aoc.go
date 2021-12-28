package aoc

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Puzzler is a function that, given a channel of line-oriented input, returns
// the answer to a puzzle, doing any descriptive logging to log.
type Puzzler func(r io.Reader, l *log.Logger) string

// Day represents a single Day of Advent of Code
type Day struct {
	number       int
	defaultInput string
	puzzles      []Puzzler
}

// Year represents a single year of AOC
type Year struct {
	number int
	days   []Day
}

// Day(int) returns the given day + a flag indicating whether it was found
func (y *Year) Day(number int) (Day, bool) {
	for _, d := range y.days {
		if d.number == number {
			return d, true
		}
	}
	return Day{}, false
}

func (y *Year) String() string {
	return fmt.Sprintf("%d", y.number)
}

func (d *Day) DefaultInput() string {
	return d.defaultInput
}

func (d *Day) Puzzles() []Puzzler {
	return d.puzzles
}

func (d *Day) String() string {
	return fmt.Sprintf("%d", d.number)
}

func NewDay(number int, defaultInput string, puzzles ...Puzzler) Day {
	return Day{number, defaultInput, puzzles}
}

func NewYear(number int, days ...Day) Year {
	return Year{number, days}
}

func Run(p Puzzler, input io.Reader) {

	l := log.New(os.Stderr, "", log.Default().Flags())

	result := p(input, l)

	fmt.Println(result)
}
