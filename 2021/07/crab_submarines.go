package d07

import (
	"bufio"
	_ "embed"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(7, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	positions := parseInput(r)

	_, lowestCost := solve(positions, getNaiveCostToMoveToPosition)

	return strconv.Itoa(lowestCost)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	positions := parseInput(r)

	_, lowestCost := solve(positions, getCostToMoveToPosition)

	return strconv.Itoa(lowestCost)
}

func solve(positions []int, costFunc func(int, []int) int) (int, int) {
	var min, max int
	for _, x := range positions {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}

	costsByPosition := make(map[int]int)
	allPositions := make([]int, 0, max-min)

	for x := min; x <= max; x++ {
		costsByPosition[x] = costFunc(x, positions)
		allPositions = append(allPositions, x)
	}

	sort.Slice(allPositions, func(i, j int) bool { return costsByPosition[allPositions[i]] < costsByPosition[allPositions[j]] })

	bestPosition := allPositions[0]
	lowestCost := costsByPosition[bestPosition]

	return bestPosition, lowestCost

}

func getNaiveCostToMoveToPosition(targetX int, positions []int) int {
	var result int

	for _, x := range positions {
		cost := int(math.Abs(float64(x - targetX)))
		result += cost
	}

	return result
}

func getCostToMoveToPosition(targetX int, positions []int) int {
	var result int

	for _, x := range positions {
		cost := costOfMove(x, targetX)
		result += cost
	}

	return result
}

func costOfMove(from, to int) int {
	distance := int(math.Abs(float64(from - to)))
	var result int

	for i := 1; i <= distance; i++ {
		result += i
	}

	return result
}

func parseInput(r io.Reader) []int {
	var result []int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		tokens := strings.Split(line, ",")
		for _, token := range tokens {
			value, err := strconv.ParseInt(token, 10, 32)
			if err != nil {
				continue
			}
			result = append(result, int(value))
		}
	}

	return result
}
