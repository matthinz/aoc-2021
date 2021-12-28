package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
	y2021 "github.com/matthinz/aoc-golang/2021"
)

const FirstYear = 2015
const MaxDays = 25

var AllYears = map[int]func() aoc.Year{
	2021: y2021.New,
}

func main() {

	yearNumbers, dayNumbers, err := parseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}

	for _, yearNumber := range yearNumbers {
		factory, found := AllYears[yearNumber]
		if !found {
			panic(fmt.Sprintf("Invalid year: %d", yearNumber))
		}
		year := factory()

		for _, dayNumber := range dayNumbers {
			day, found := year.Day(dayNumber)
			if !found {
				panic(fmt.Sprintf("Year %d, day %d not found", yearNumber, dayNumber))
			}

			runDay(&day)
		}
	}
}

func parseArgs(args []string) ([]int, []int, error) {

	var years []int
	var days []int

	for _, arg := range args {
		num, err := strconv.ParseInt(arg, 10, 16)
		if err != nil {
			return years, days, fmt.Errorf("Invalid argument: %s", arg)
		}

		if num < 1 {
			return years, days, fmt.Errorf("Invalid argument: %s", arg)
		}

		if num >= FirstYear {
			years = append(years, int(num))
		} else if num <= MaxDays {
			days = append(days, int(num))
		}
	}

	if len(years) == 0 {
		for y := range AllYears {
			years = append(years, y)
		}
	}

	if len(days) == 0 {
		for d := 1; d <= MaxDays; d++ {
			days = append(days, d)
		}
	}

	return years, days, nil
}

func runDay(day *aoc.Day) {
	var input io.ReadSeeker

	stat, err := os.Stdin.Stat()
	if err == nil {
		isTTY := (stat.Mode() & os.ModeCharDevice) != 0
		if isTTY {
			// no file, use the default
			input = strings.NewReader(day.DefaultInput())
		}
	}

	if input == nil {
		input = os.Stdin
	}

	for _, p := range day.Puzzles() {
		input.Seek(0, io.SeekStart)
		aoc.Run(p, input)
	}

}
