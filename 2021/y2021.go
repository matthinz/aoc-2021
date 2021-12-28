package y2021

import (
	aoc "github.com/matthinz/aoc-golang"
	d01 "github.com/matthinz/aoc-golang/2021/01"
	d02 "github.com/matthinz/aoc-golang/2021/02"
	d03 "github.com/matthinz/aoc-golang/2021/03"
	d13 "github.com/matthinz/aoc-golang/2021/13"
)

func New() aoc.Year {
	return aoc.NewYear(
		2021,
		d01.New(),
		d02.New(),
		d03.New(),
		d13.New(),
	)
}
