package d20

import (
	"fmt"
	"strings"
	"testing"
)

func TestChunkToIndex(t *testing.T) {
	chunk := [3][3]bool{
		{false, false, false},
		{true, false, false},
		{false, true, false},
	}
	expected := 34
	actual := chunkToIndex(chunk)
	if actual != expected {
		t.Errorf("Expected %d (%b), got %d (%b)", expected, expected, actual, actual)
	}
}

func TestEnhance(t *testing.T) {
	input := strings.TrimSpace(`
..#.#..#####.#.#.#.###.##.....###.##.#..###.####..#####..#....#..#..##..###..######.###...####..#..#####..##..#.#####...##.#.#..#.##..#.#......#.###.######.###.####...#.##.##..#..#..#####.....#.#....###..#.##......#.....#..#..#..##..#...##.######.####.####.#.#...#.......#..#.#.#...####.##.#......#..#...##.#.##..#...##.#.##..###.#......#.#.......#.#.#.####.###.##...#.....####.#..#..#.##.#....##..#.####....##...##..#...#......#.#.......#.......##..####..#...#.#.#...##..#.#..###..#####........#..####......#..#

#..#.
#....
##..#
..#..
..###
	`)

	img, algorithm := parseInput(strings.NewReader(input))

	enhanced := enhance(&img, algorithm)

	for y := 0; y < enhanced.height; y++ {
		fmt.Println()
		for x := 0; x < enhanced.width; x++ {
			if enhanced.pixels[y][x] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
	}

	expected := 24
	actual := countLitPixels(enhanced)

	if actual != expected {
		t.Errorf("Wrong # of lit pixels. Expected %d, got %d", expected, actual)
	}

}
