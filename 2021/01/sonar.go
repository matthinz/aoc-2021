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

	s := bufio.NewScanner(r)

	var prevValue *int
	var increases int

	for s.Scan() {
		valueAsInt64, err := strconv.ParseInt(s.Text(), 10, 32)
		if err != nil {
			continue
		}

		value := int(valueAsInt64)

		if prevValue != nil {

			increased := value > *prevValue
			decreased := value < *prevValue

			if increased {
				increases++
				log.Printf("%d: increased\n", valueAsInt64)
			} else if decreased {
				log.Printf("%d: decreased\n", valueAsInt64)
			}

		}

		prevValue = &value
	}

	return strconv.Itoa(increases)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	const WindowSize = 3

	scanner := bufio.NewScanner(r)
	var window []int
	var prevSum *int
	var increases int

	for scanner.Scan() {
		valueAsInt64, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			continue
		}

		value := int(valueAsInt64)

		window = append(window, value)

		if len(window) == WindowSize {

			sum := sumValues(window)

			if prevSum != nil {
				increased := sum > *prevSum
				decreased := sum < *prevSum

				if increased {
					log.Printf("%s = %d (increased)\n", formatEquation(window), sum)
					increases++
				} else if decreased {
					log.Printf("%s = %d (decreased)\n", formatEquation(window), sum)
				}
			} else {
				log.Printf("%s = %d (n/a)\n", formatEquation(window), sum)
			}

			prevSum = &sum

			window = window[1:]
		}
	}

	return strconv.Itoa(increases)

}

func sumValues(values []int) int {
	var result int
	for _, v := range values {
		result += v
	}
	return result
}

func formatEquation(values []int) string {
	builder := strings.Builder{}
	for index, v := range values {
		if index > 0 {
			builder.WriteString("+")
		}
		builder.WriteString(strconv.Itoa(v))
	}
	return builder.String()
}
