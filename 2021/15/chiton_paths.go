package d15

import (
	"bufio"
	_ "embed"
	"io"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/matthinz/aoc-golang"
)

type point struct {
	x         int
	y         int
	risk      int
	totalRisk int
}

const unknownRisk = -1

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(15, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	grid := parseInput(r)
	lowestTotalRisk := solveDijkstra(grid, l)
	return strconv.Itoa(lowestTotalRisk)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	grid := parseInput(r)
	grid = inflateGrid(grid, 5)
	lowestTotalRisk := solveDijkstra(grid, l)
	return strconv.Itoa(lowestTotalRisk)
}

func inflateGrid(grid *[][]int, inflationFactor int) *[][]int {
	height := len(*grid)
	width := len((*grid)[0])
	inflatedGrid := make([][]int, height*inflationFactor)

	for y := 0; y < height; y++ {
		for i := 0; i < inflationFactor; i++ {
			destY := (height * i) + y
			inflatedGrid[destY] = make([]int, width*inflationFactor)

			for x := 0; x < width; x++ {
				for j := 0; j < inflationFactor; j++ {

					destX := (width * j) + x

					value := ((*grid)[y][x] + i + j)
					if value > 9 {
						value = value % 9
					}

					inflatedGrid[destY][destX] = value
				}
			}
		}
	}

	return &inflatedGrid
}

func solveDijkstra(grid *[][]int, l *log.Logger) int {

	height := len(*grid)
	width := len((*grid)[0])

	// toVisit is a list of points, sorted by totalRisk -- from
	// highest risk to lowest risk
	toVisit := make([]point, 0, width*height)

	// going in reverse here so that toVisit ends up with
	// lowest totalRisk at the end
	for y := height - 1; y >= 0; y-- {
		for x := width - 1; x >= 0; x-- {
			totalRisk := unknownRisk

			if x == 0 && y == 0 {
				// starting location = no risk
				totalRisk = 0
			}

			p := point{
				x:         x,
				y:         y,
				risk:      (*grid)[y][x],
				totalRisk: totalRisk,
			}

			toVisit = append(toVisit, p)
		}
	}

	visited := make([]point, 0, width*height)

	start := time.Now()

	for len(toVisit) > 0 {

		if rand.Float64() < 0.001 {
			left := len(toVisit)
			duration := time.Now().Sub(start)
			processed := (width * height) - left
			secondsToProcessOne := duration.Seconds() / float64(processed)
			remainingSeconds := secondsToProcessOne * float64(left)

			if left > 10000 {
				l.Printf("%d to visit, %v elapsed, ~%fs remain\n", left, duration, remainingSeconds)
			}
		}

		// Get the last point in toVisit, which will be the *least* risky
		point, nextToVisit := popPoint(toVisit)
		toVisit = nextToVisit

		visited = append(visited, point)

		needToSort := false

		for i := range toVisit {

			isNorthNeighbor := (toVisit[i].x == point.x && toVisit[i].y == point.y-1)
			isSouthNeighbor := (toVisit[i].x == point.x && toVisit[i].y == point.y+1)
			isEastNeighbor := (toVisit[i].x == point.x+1 && toVisit[i].y == point.y)
			isWestNeighbor := (toVisit[i].x == point.x-1 && toVisit[i].y == point.y)

			isNeighbor := isNorthNeighbor || isSouthNeighbor || isEastNeighbor || isWestNeighbor

			if !isNeighbor {
				continue
			}

			// we've found an unvisited neighbor compute its total risk, then update
			// its position in <toVisit>
			neighbor := &toVisit[i]
			newRisk := point.totalRisk + neighbor.risk
			if neighbor.totalRisk == unknownRisk || neighbor.totalRisk > newRisk {
				neighbor.totalRisk = newRisk
				needToSort = true
			}
		}

		if !needToSort {
			continue
		}

		// keep toVisit sorted such that the least risky points are at the end
		sort.Slice(toVisit, func(i, j int) bool {

			if toVisit[i].totalRisk == toVisit[j].totalRisk {
				return false
			}

			// unknown risk = sort at the beginning
			if toVisit[i].totalRisk == unknownRisk {
				return true
			}

			if toVisit[j].totalRisk == unknownRisk {
				return false
			}

			// known risk = the lower the number, the higher it is sorted
			return toVisit[i].totalRisk > toVisit[j].totalRisk
		})
	}

	// visited now contains all points, each with the lowest possible risk
	// to get to them
	for _, p := range visited {
		if p.x == width-1 && p.y == height-1 {
			return p.totalRisk
		}
	}

	panic("Could not determine total risk")
}

func popPoint(slice []point) (point, []point) {
	sliceLen := len(slice)
	p := slice[sliceLen-1]
	newSlice := make([]point, sliceLen-1)
	copy(newSlice, slice)
	return p, newSlice
}

func parseInput(r io.Reader) *[][]int {
	var grid [][]int
	var width int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}
		row := make([]int, 0, width)
		for _, r := range line {
			num, err := strconv.ParseInt(string(r), 10, 8)
			if err != nil {
				panic(err)
			}
			row = append(row, int(num))
		}
		grid = append(grid, row)
	}

	return &grid
}
