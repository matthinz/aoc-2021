package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x      int
	y      int
	height int
}

func main() {
	input := readInput(os.Stdin)
	lowPoints := getLowPoints(input)

	sumOfRiskLevels := 0
	for _, p := range lowPoints {
		sumOfRiskLevels += p.height + 1
	}

	fmt.Printf("Sum of risk levels: %d\n", sumOfRiskLevels)

	basins := findBasins(input)

	sort.Slice(basins, func(i, j int) bool {
		return len(basins[i]) > len(basins[j])
	})

	sizes := 1
	for i, b := range basins[0:3] {
		fmt.Printf("%d. %d\n", i+1, len(b))
		sizes *= len(b)
	}

	fmt.Printf("3 largest multiplied: %d\n", sizes)

}

func findBasins(input [][]int) [][]point {
	var result [][]point

	for _, p := range getLowPoints(input) {
		basin := buildBasin(input, p, make([]point, 0))
		result = append(result, basin)
	}

	return result
}

func havePointsInCommon(basin1 []point, basin2 []point) bool {
	for _, a := range basin1 {
		if hasPoint(basin2, a.x, a.y) {
			return true
		}
	}
	return false
}

func hasPoint(slice []point, x, y int) bool {
	for _, p := range slice {
		if p.x == x && p.y == y {
			return true
		}
	}
	return false
}

func buildBasin(input [][]int, p point, basin []point) []point {

	if p.height >= 9 {
		return basin
	}

	for _, b := range basin {
		if b.x == p.x && b.y == p.y {
			// this basin already contains this point
			return basin
		}
	}

	result := append(basin, p)

	// above
	if p.y > 0 {
		result = buildBasin(input, point{p.x, p.y - 1, input[p.y-1][p.x]}, result)
	}

	// below
	if p.y < len(input)-1 {
		result = buildBasin(input, point{p.x, p.y + 1, input[p.y+1][p.x]}, result)
	}

	// left
	if p.x > 0 {
		result = buildBasin(input, point{p.x - 1, p.y, input[p.y][p.x-1]}, result)
	}

	// right
	if p.x < len(input[p.y])-1 {
		result = buildBasin(input, point{p.x + 1, p.y, input[p.y][p.x+1]}, result)
	}

	return result
}

func getLowPoints(input [][]int) []point {
	var result []point

	for y, row := range input {
		for x, height := range row {

			adjacents := make([]int, 0, 4)

			// above
			if y > 0 {
				adjacents = append(adjacents, input[y-1][x])
			}

			// below
			if y < len(input)-1 {
				adjacents = append(adjacents, input[y+1][x])
			}

			// left
			if x > 0 {
				adjacents = append(adjacents, input[y][x-1])
			}

			// right
			if x < len(row)-1 {
				adjacents = append(adjacents, input[y][x+1])
			}

			isLow := true

			for _, adjacent := range adjacents {
				if height >= adjacent {
					isLow = false
				}
			}

			if isLow {
				result = append(result, point{x, y, height})
			}
		}
	}

	return result
}

func readInput(r io.Reader) [][]int {
	var result [][]int
	var scanner = bufio.NewScanner(r)

	for scanner.Scan() {
		lineOfText := strings.TrimSpace(scanner.Text())
		row := make([]int, 0, len(lineOfText))

		for _, r := range lineOfText {
			num, err := strconv.ParseInt(string(r), 10, 8)
			if err != nil {
				continue
			}
			row = append(row, int(num))
		}
		result = append(result, row)
	}

	return result
}
