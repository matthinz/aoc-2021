package y2021

import (
	aoc "github.com/matthinz/aoc-golang"
	d01 "github.com/matthinz/aoc-golang/2021/01"
	d02 "github.com/matthinz/aoc-golang/2021/02"
	d03 "github.com/matthinz/aoc-golang/2021/03"
	d04 "github.com/matthinz/aoc-golang/2021/04"
	d05 "github.com/matthinz/aoc-golang/2021/05"
	d06 "github.com/matthinz/aoc-golang/2021/06"
	d07 "github.com/matthinz/aoc-golang/2021/07"
	d08 "github.com/matthinz/aoc-golang/2021/08"
	d09 "github.com/matthinz/aoc-golang/2021/09"
	d10 "github.com/matthinz/aoc-golang/2021/10"
	d11 "github.com/matthinz/aoc-golang/2021/11"
	d12 "github.com/matthinz/aoc-golang/2021/12"
	d13 "github.com/matthinz/aoc-golang/2021/13"
)

func New() aoc.Year {
	return aoc.NewYear(
		2021,
		d01.New(),
		d02.New(),
		d03.New(),
		d04.New(),
		d05.New(),
		d06.New(),
		d07.New(),
		d08.New(),
		d09.New(),
		d10.New(),
		d11.New(),
		d12.New(),
		d13.New(),
	)
}
