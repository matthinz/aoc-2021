package d02

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
	return aoc.NewDay(2, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	var x int
	var y int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			continue
		}
		command := tokens[0]
		value, err := strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			panic(err)
		}

		switch command {
		case "forward":
			x += int(value)
			break
		case "down":
			y += int(value)
			break
		case "up":
			y -= int(value)
			break
		}

		l.Printf("%s %d: %d x %d = %d\n", command, value, x, y, x*y)
	}
	return strconv.Itoa(x * y)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	var x int
	var y int
	var aim int

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			continue
		}
		command := tokens[0]
		value, err := strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			panic(err)
		}

		switch command {
		case "forward":
			x += int(value)
			y += aim * int(value)
			break
		case "down":
			aim += int(value)
			break
		case "up":
			aim -= int(value)
			break
		}
	}
	return strconv.Itoa(x * y)
}
