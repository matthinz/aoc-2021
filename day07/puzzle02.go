package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	positions := parseInput(os.Stdin)

	var min, max int
	for _, x := range positions {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}

	fmt.Printf("min: %d, max: %d\n", min, max)

	costsByPosition := make(map[int]uint)
	allPositions := make([]int, 0, len(positions))

	// There may be intermediate positions not included in the input -- we want
	// to make sure we consider those as places to go
	for x := min; x <= max; x++ {
		_, exists := costsByPosition[x]
		if !exists {
			costsByPosition[x] = getCostToMoveToPosition(x, positions)
			allPositions = append(allPositions, x)
		}
	}

	sort.Slice(allPositions, func(i, j int) bool { return costsByPosition[allPositions[j]] < costsByPosition[allPositions[i]] })

	for _, x := range allPositions {
		fmt.Printf("%d: %d\n", x, costsByPosition[x])
	}

}

func getCostToMoveToPosition(targetX int, positions []int) uint {
	var result uint

	for _, x := range positions {
		cost := costOfMove(x, targetX)
		result += cost
	}

	return result
}

func costOfMove(from, to int) uint {
	distance := int(math.Abs(float64(from - to)))
	var result uint

	for i := 1; i <= distance; i++ {
		result += uint(i)
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
