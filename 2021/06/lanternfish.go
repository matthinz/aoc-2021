package d06

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
	return aoc.NewDay(6, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	numbers := parseNumbers(r)
	return strconv.Itoa(simulate(numbers, 80))
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	numbers := parseNumbers(r)
	return strconv.Itoa(simulate(numbers, 256))
}

func simulate(numbers []int, days int) int {

	const MaxTimer = 8
	timers := make([]int, MaxTimer+1)

	for _, timer := range numbers {
		timers[timer]++
	}

	for day := 1; day <= days; day++ {

		nextTick := make([]int, MaxTimer+1)

		for i := 0; i <= MaxTimer; i++ {
			count := timers[i]
			switch i {
			case 0:
				nextTick[6] += count
				nextTick[8] += count
				break
			default:
				nextTick[i-1] += count
				break
			}
		}

		timers = nextTick
	}

	return sum(timers)
}

func sum(values []int) int {
	var result int
	for _, value := range values {
		result += value
	}
	return result
}

func parseNumbers(r io.Reader) []int {

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		var result []int
		for _, token := range strings.Split(line, ",") {
			value, err := strconv.ParseInt(token, 10, 32)
			if err != nil {
				continue
			}
			result = append(result, int(value))
		}
		return result
	}

	return []int{}
}
