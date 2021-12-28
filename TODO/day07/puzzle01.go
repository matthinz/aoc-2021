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

	costsByPosition := make(map[int]int)

	for x := min; x <= max; x++ {
		costsByPosition[x] = getCostToMoveToPosition(x, positions)
	}

	sort.Slice(positions, func(i, j int) bool { return costsByPosition[positions[i]] < costsByPosition[positions[j]] })

	for _, x := range positions {
		fmt.Printf("%d: %d\n", x, costsByPosition[x])
	}

}

func getCostToMoveToPosition(targetX int, positions []int) int {
	var result int

	for _, x := range positions {
		cost := int(math.Abs(float64(x - targetX)))
		result += cost
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
