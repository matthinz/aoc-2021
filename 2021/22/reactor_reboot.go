package d22

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type point struct {
	x, y, z int
}

type cuboid struct {
	position point
	size     point
}

type step struct {
	cuboid cuboid
	turnOn bool
}

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(22, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	steps := parseInput(r)

	reactor := applyStepsUsingBruteForce(steps)

	ct := countCubesOnUsingBruteForce(reactor)

	return strconv.FormatUint(uint64(ct), 10)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// Brute force solution

func applyStepsUsingBruteForce(steps []step) *[][][]bool {
	var reactor *[][][]bool
	for _, s := range steps {
		reactor = applyStepUsingBruteForce(s, reactor)
	}
	return reactor
}

func applyStepUsingBruteForce(s step, reactor *[][][]bool) *[][][]bool {

	const minX = -50
	const maxX = 50
	const minY = -50
	const maxY = 50
	const minZ = -50
	const maxZ = 50

	var workingReactor [][][]bool

	if reactor == nil {
		width := maxX - minX + 1
		workingReactor = make([][][]bool, width)
		for x := 0; x < width; x++ {
			height := maxY - minY + 1
			workingReactor[x] = make([][]bool, height)
			for y := 0; y < height; y++ {
				depth := maxZ - minZ + 1
				workingReactor[x][y] = make([]bool, depth)
			}
		}
	} else {
		workingReactor = *reactor
	}

	position := s.cuboid.position
	size := s.cuboid.size

	if position.x < minX || position.x > maxX || position.x+size.x > maxX {
		return &workingReactor
	}

	if position.y < minY || position.y > maxY || position.y+size.y > maxY {
		return &workingReactor
	}

	if position.z < minZ || position.z > maxZ || position.z+size.z > maxZ {
		return &workingReactor
	}

	for x := position.x - minX; x < position.x+size.x-minX; x++ {
		for y := position.y - minY; y < position.y+size.y-minY; y++ {
			for z := position.z - minZ; z < position.z+size.z-minZ; z++ {
				workingReactor[x][y][z] = s.turnOn
			}
		}
	}

	return &workingReactor
}

func countCubesOnUsingBruteForce(reactor *[][][]bool) uint {
	var count uint
	for x := 0; x < len(*reactor); x++ {
		for y := 0; y < len((*reactor)[x]); y++ {
			for z := 0; z < len((*reactor)[x][y]); z++ {
				if (*reactor)[x][y][z] {
					count++
				}
			}
		}
	}
	return count
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) []step {
	var steps []step

	rangeRx := "(-?\\d+)\\.\\.(-?\\d+)"
	rx := regexp.MustCompile(
		fmt.Sprintf(
			"(on|off) x=%s,y=%s,z=%s",
			rangeRx, rangeRx, rangeRx,
		),
	)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		m := rx.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		c, err := parseCuboid(m[2], m[3], m[4], m[5], m[6], m[7])
		if err != nil {
			panic(err)
		}

		step := step{
			cuboid: c,
			turnOn: m[1] == "on",
		}

		steps = append(steps, step)
	}

	return steps
}

func parseCuboid(x1, x2, y1, y2, z1, z2 string) (cuboid, error) {

	inputs := [6]string{x1, x2, y1, y2, z1, z2}
	values := [6]int{}

	for i := range inputs {
		value, err := strconv.ParseInt(inputs[i], 10, 32)
		if err != nil {
			return cuboid{}, err
		}

		values[i] = int(value)
	}

	c := cuboid{
		position: point{
			values[0], values[2], values[4],
		},
		size: point{
			values[1] - values[0] + 1,
			values[3] - values[2] + 1,
			values[5] - values[4] + 1,
		},
	}

	return c, nil
}
