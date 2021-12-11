package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	input := parseInput(os.Stdin)
	stepIndex := 0
	totalFlashes := 0

	height := len(input)
	width := len(input[0])
	area := width * height

	for {
		stepIndex++

		nextInput, flashes := step(input)
		totalFlashes += flashes

		fmt.Printf("%d: %d (%d)\n", stepIndex, flashes, totalFlashes)

		if flashes == area {
			fmt.Printf("All %d octopuses flashed!!\n", area)
			break
		}

		input = nextInput
	}
}

func parseInput(r io.Reader) [][]int {
	var result [][]int
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}
		row := make([]int, len(line))
		for i, r := range line {
			num, err := strconv.ParseInt(string(r), 10, 8)
			if err != nil {
				panic(err)
			}
			row[i] = int(num)
		}
		result = append(result, row)
	}
	return result
}

func step(input [][]int) ([][]int, int) {
	result := increment(input)

	totalFlashes := 0
	for {
		next, flashes := flash(result)
		result = next
		totalFlashes += flashes
		if flashes == 0 {
			break
		}
	}

	return reset(result), totalFlashes
}

func reset(input [][]int) [][]int {
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			if input[y][x] > 9 || input[y][x] == -1 {
				input[y][x] = 0
			}
		}
	}
	return input
}

func increment(input [][]int) [][]int {
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			input[y][x]++
		}
	}
	return input
}

func flash(input [][]int) ([][]int, int) {

	flashes := 0

	height := len(input)
	width := len(input[0])

	// make a copy of the input for us to modify
	// TODO: will go's copy() function work here?
	result := make([][]int, height)
	for y := 0; y < height; y++ {
		if len(input[y]) != width {
			panic("row has wrong width")
		}
		result[y] = make([]int, width)
		for x := 0; x < len(input[y]); x++ {
			result[y][x] = input[y][x]
		}
	}

	// flash everything in the input that is > 9
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if input[y][x] <= 9 {
				continue
			}

			for deltaY := -1; deltaY <= 1; deltaY++ {
				for deltaX := -1; deltaX <= 1; deltaX++ {

					yValid := y+deltaY >= 0 && y+deltaY < height
					xValid := x+deltaX >= 0 && x+deltaX < width

					if yValid && xValid {

						hasAlreadyFlashed := result[y+deltaY][x+deltaX] == -1

						if hasAlreadyFlashed {
							// this octopus has already flashed on this tick
							continue
						}

						result[y+deltaY][x+deltaX]++
					}

				}
			}

			result[y][x] = -1 // sentinel value indicating "don't mess with this one again"

			flashes++
		}
	}

	return result, flashes
}

func printGrid(grid [][]int) {
	for y := 0; y < len(grid); y++ {
		fmt.Println()
		for x := 0; x < len(grid[y]); x++ {
			fmt.Print(grid[y][x])
		}
	}
	fmt.Println()
}
