package d22

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
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
	on       bool
}

type interval struct {
	start, end int

	// Sorted set of indices of cuboids that correspond to this interval
	cuboidIndices []int
}

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(22, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	cuboids := parseInput(r)

	reactor := initializeReactorUsingBruteForce(cuboids)

	ct := countCubesOnUsingBruteForce(reactor)

	return strconv.FormatUint(uint64(ct), 10)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	cuboids := parseInput(r)

	normalized := initializeReactor(cuboids)

	var ct uint

	for _, c := range normalized {
		ct += (uint(c.size.x) * uint(c.size.y) * uint(c.size.z))
	}

	return strconv.FormatUint(uint64(ct), 10)
}

////////////////////////////////////////////////////////////////////////////////
// Non-brute force solution

// takes a set of cuboids and returns a normalized set of non-overlapping
// cuboids that have been turned on
func initializeReactor(cuboids []cuboid) []cuboid {
	return []cuboid{}
}

// given a set of cuboids along with functions to pull off specific coordinate
// + extent values (e.g. x and width, y and height, etc.), returns the set of
// intervals generated.
func buildIntervals(cuboids []cuboid, getCoordinate func(c cuboid) int, getExtent func(c cuboid) int) []interval {
	// // these maps connect values to the indices of <steps> that contain those values
	// values, valuesToStepIndices := buildValueMap(cuboids, func(c cuboid) []int {
	// 	return []int{
	// 		getCoordinate(c),
	// 		getCoordinate(c) + getExtent(c) - 1,
	// 	}
	// })

	result := []interval{}
	return result
}

// given a set of cuboids and a function to read coordinate values them, returns
// a sorted slice of those values and a map connecting those values to the
// cuboid indices that contained those values
func buildValueMap(cuboids []cuboid, f func(c cuboid) []int) ([]int, map[int][]int) {
	values := make([]int, 0)
	valueToIndexMap := make(map[int][]int)

	for cuboidIndex, c := range cuboids {
		for _, value := range f(c) {
			valueToIndexMap[value] = append(valueToIndexMap[value], cuboidIndex)
			values = append(values, value)
		}
	}

	sort.Ints(values)

	return values, valueToIndexMap
}

////////////////////////////////////////////////////////////////////////////////
// Brute force solution

func initializeReactorUsingBruteForce(cuboids []cuboid) *[][][]bool {
	var reactor *[][][]bool
	for _, c := range cuboids {
		reactor = processCuboidUsingBruteForce(c, reactor)
	}
	return reactor
}

func processCuboidUsingBruteForce(c cuboid, reactor *[][][]bool) *[][][]bool {

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

	if c.position.x < minX || c.position.x > maxX || c.position.x+c.size.x > maxX {
		return &workingReactor
	}

	if c.position.y < minY || c.position.y > maxY || c.position.y+c.size.y > maxY {
		return &workingReactor
	}

	if c.position.z < minZ || c.position.z > maxZ || c.position.z+c.size.z > maxZ {
		return &workingReactor
	}

	for x := c.position.x - minX; x < c.position.x+c.size.x-minX; x++ {
		for y := c.position.y - minY; y < c.position.y+c.size.y-minY; y++ {
			for z := c.position.z - minZ; z < c.position.z+c.size.z-minZ; z++ {
				workingReactor[x][y][z] = c.on
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

func parseInput(r io.Reader) []cuboid {
	var cuboids []cuboid

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

		c.on = m[1] == "on"

		cuboids = append(cuboids, c)
	}

	return cuboids
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
