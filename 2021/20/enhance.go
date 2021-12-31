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
	width         int
	height        int
	pixels        [][]bool
	infinitePixel bool
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
	img, algorithm := parseInput(r)

	enhanced := &img

	for i := 0; i < 50; i++ {
		enhanced = enhance(enhanced, algorithm)
	}

	count := countLitPixels(enhanced)

	return strconv.Itoa(count)
}

////////////////////////////////////////////////////////////////////////////////
// enhance!

func enhance(img *image, algorithm []bool) *image {

	result := image{
		width:  img.width + 2,
		height: img.height + 2,
	}
	result.pixels = make([][]bool, result.height)

	if img.infinitePixel {
		// new infinite pixel will be algorithm at 0x111111111 (511)
		result.infinitePixel = algorithm[511]
	} else {
		// new infinite pixel will be algorithm at 0
		result.infinitePixel = algorithm[0]
	}

	for destY := 0; destY < result.height; destY++ {
		srcY := destY - 1
		result.pixels[destY] = make([]bool, result.width)

		for destX := 0; destX < result.width; destX++ {
			srcX := destX - 1
			chunk := extractChunk(img, srcX, srcY)

			index := chunkToIndex(chunk)

			result.pixels[destY][destX] = algorithm[index]
		}
	}

	return &result
}

func trimImage(img *image) *image {

	leftMost := -1
	topMost := -1
	rightMost := -1
	bottomMost := -1

	for y := 0; y < img.height; y++ {
		for x := 0; x < img.width; x++ {

			pixel := img.pixels[y][x]
			if !pixel {
				continue
			}

			if leftMost == -1 || x < leftMost {
				leftMost = x
			}

			if rightMost == -1 || x > rightMost {
				rightMost = x
			}

			if topMost == -1 || y < topMost {
				topMost = y
			}

			if bottomMost == -1 || y > bottomMost {
				bottomMost = y
			}
		}
	}

	newWidth := rightMost - leftMost + 1
	newHeight := bottomMost - topMost + 1

	if newWidth < 0 || newHeight < 0 {
		return &image{
			width:  0,
			height: 0,
		}
	}

	result := image{
		width:  newWidth,
		height: newHeight,
	}

	result.pixels = make([][]bool, result.height)

	for srcY := topMost; srcY <= bottomMost; srcY++ {
		destY := srcY - topMost
		result.pixels[destY] = make([]bool, result.width)

		for srcX := leftMost; srcX <= rightMost; srcX++ {
			destX := srcX - leftMost
			result.pixels[destY][destX] = img.pixels[srcY][srcX]
		}
	}

	return &result
}

func chunkToIndex(chunk [][]bool) int {
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
func extractChunk(img *image, x int, y int) [][]bool {

	chunk := make([][]bool, 3)

	for dY := -1; dY <= 1; dY++ {
		chunkY := dY + 1
		chunk[chunkY] = make([]bool, 3)

		for dX := -1; dX <= 1; dX++ {

			chunkX := dX + 1

			actualY := y + dY
			actualX := x + dX

			var pixel bool

			if actualY < 0 || actualY >= img.height {
				pixel = img.infinitePixel
			} else if actualX < 0 || actualX >= img.width {
				pixel = img.infinitePixel
			} else {
				pixel = img.pixels[actualY][actualX]
			}

			chunk[chunkY][chunkX] = pixel

		}
	}

	return chunk
}

func countLitPixels(img *image) int {
	if img.infinitePixel {
		panic("infinite pixels are lit")
	}

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

func printPixels(pixels *[][]bool) string {
	sb := strings.Builder{}
	for y := 0; y < len(*pixels); y++ {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		for x := 0; x < len((*pixels)[y]); x++ {
			pixel := (*pixels)[y][x]
			if pixel {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
	}
	return sb.String()
}
