package d01

import (
	"bufio"
	_ "embed"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(1, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	numbers := parseInput(r)

	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			if numbers[i]+numbers[j] == 2020 {
				return strconv.Itoa(numbers[i] * numbers[j])
			}
		}
	}

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	numbers := parseInput(r)

	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			for k := i + 1; k < len(numbers); k++ {
				if numbers[i]+numbers[j]+numbers[k] == 2020 {
					return strconv.Itoa(numbers[i] * numbers[j] * numbers[k])
				}
			}

		}
	}

	return ""
}

func parseInput(r io.Reader) []int {
	result := make([]int, 0)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}
		parsed, err := strconv.ParseInt(line, 10, 32)
		if err != nil {
			panic(err)
		}
		result = append(result, int(parsed))
	}

	return result
}
