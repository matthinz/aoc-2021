package d03

import (
	"bufio"
	_ "embed"
	"io"
	"log"
	"strconv"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(3, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	numbers := parseInput(r)

	gamma, epsilon := calculateGammaAndEpsilon(numbers)

	return strconv.Itoa(gamma * epsilon)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	numbers := parseInput(r)

	o2GeneratorRating := FindO2GeneratorRating(numbers)
	co2ScrubberRating := FindCo2ScrubberRating(numbers)

	return strconv.Itoa(o2GeneratorRating * co2ScrubberRating)

}

func parseInput(r io.Reader) []int {
	scanner := bufio.NewScanner(r)
	var result []int

	for scanner.Scan() {
		token := scanner.Text()
		if len(token) == 0 {
			continue
		}
		value, err := strconv.ParseInt(token, 2, 64)
		if err != nil {
			panic(err)
		}
		result = append(result, int(value))
	}
	return result
}
