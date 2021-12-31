package d20

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed input
var actualRealInput string

func TestChunkToIndex(t *testing.T) {
	chunk := [][]bool{
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

func TestExtractChunk(t *testing.T) {

	img := image{
		width:  3,
		height: 3,
		pixels: [][]bool{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}}

	chunk := extractChunk(&img, -1, -1)

	expected := strings.TrimSpace(`
...
...
..#
		`)

	actual := printPixels(&chunk)

	if actual != expected {
		t.Log(actual)
		t.Error("Extract failed")
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

	expected := 24
	actual := countLitPixels(enhanced)

	if actual != expected {
		t.Log(printPixels(&enhanced.pixels))
		t.Errorf("Wrong # of lit pixels. Expected %d, got %d", expected, actual)
	}

	enhanced = enhance(enhanced, algorithm)

	expected = 35
	actual = countLitPixels(enhanced)

	if actual != expected {
		t.Log(printPixels(&enhanced.pixels))
		t.Errorf("Wrong # of lit pixels after 2nd enhance. Expected %d, got %d", expected, actual)
	}

}

func TestEnhanceWithRealInput(t *testing.T) {

	img, algorithm := parseInput(strings.NewReader(actualRealInput))

	enhanced := enhance(&img, algorithm)
	enhanced = enhance(enhanced, algorithm)

	actual := countLitPixels(enhanced)

	wrongAnswers := []int{
		6224,
		5687,
		5602,
	}

	for _, wrongAnswer := range wrongAnswers {
		if actual == wrongAnswer {
			t.Errorf("The answer is _not_ %d", wrongAnswer)
		}
	}
}
