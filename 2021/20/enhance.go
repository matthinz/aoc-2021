package d20

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type image struct {
	width  int
	height int
	pixels [][]bool
}

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(20, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	img, algorithm := parseInput(r)

	enhanced := enhance(&img, algorithm)

	enhanced = enhance(enhanced, algorithm)

	count := countLitPixels(enhanced)

	return strconv.Itoa(count)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// enhance!

func enhance(img *image, algorithm []bool) *image {

	result := image{
		width:  img.width,
		height: img.height,
	}
	result.pixels = make([][]bool, result.height)

	for destY := 0; destY < result.height; destY++ {
		result.pixels[destY] = make([]bool, result.width)

		for destX := 0; destX < result.width; destX++ {

			srcX := destX + ((img.width - result.width) / 2)
			srcY := destY + ((img.height - result.height) / 2)

			chunk := extractChunk(img, srcX, srcY)

			index := chunkToIndex(chunk)

			result.pixels[destY][destX] = algorithm[index]
		}
	}
	return &result
}

func chunkToIndex(chunk [3][3]bool) int {
	// chunk is read left-to-right, top-to-bottom as a 9 bit binary number
	var result int

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			result <<= 1
			if chunk[y][x] {
				result = result | 1
			}
		}
	}

	return result
}

// extracts a 3x3 matrix from <img> centered on x,y
func extractChunk(img *image, x int, y int) [3][3]bool {

	chunk := [3][3]bool{}

	for dY := -1; dY <= 1; dY++ {
		for dX := -1; dX <= 1; dX++ {

			chunkY := dY + 1
			chunkX := dX + 1

			actualY := y + dY
			actualX := x + dX

			var pixel bool

			if actualY < 0 || actualY >= img.height {
				pixel = false
			} else if actualX < 0 || actualX >= img.width {
				pixel = false
			} else {
				pixel = img.pixels[actualY][actualX]
			}

			chunk[chunkY][chunkX] = pixel

		}
	}

	return chunk
}

func countLitPixels(img *image) int {
	var count int
	for y := 0; y < img.height; y++ {
		for x := 0; x < img.width; x++ {
			if img.pixels[y][x] {
				count++
			}
		}
	}
	return count
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) (image, []bool) {

	image := image{
		pixels: make([][]bool, 0),
	}
	var algorithm []bool

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())

		if len(l) == 0 {
			continue
		}

		if len(algorithm) == 0 {
			algorithm = make([]bool, len(l))
			for i, r := range l {
				switch r {
				case '#':
					algorithm[i] = true
				case '.':
					algorithm[i] = false
				default:
					panic(fmt.Sprintf("Invalid char in algorithm: %s", string(r)))
				}
			}
			continue
		}

		image.height++
		image.width = len(l)

		image.pixels = append(image.pixels, make([]bool, image.width))

		for i, r := range l {
			switch r {
			case '#':
				image.pixels[image.height-1][i] = true
			case '.':
				image.pixels[image.height-1][i] = false
			default:
				panic(fmt.Sprintf("Invalid char in image: %s", string(r)))
			}
		}
	}

	return image, algorithm
}
