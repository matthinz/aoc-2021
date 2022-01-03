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
	"time"

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
	// inclusive starting point of this interval
	start int

	// *exclusive* ending point of this interval
	end int

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

	initializationCuboids := make([]cuboid, 0)
	for _, c := range cuboids {
		if c.position.x < -50 || c.position.x > 50 {
			continue
		}
		if c.position.y < -50 || c.position.y > 50 {
			continue
		}
		if c.position.z < -50 || c.position.z > 50 {
			continue
		}
		initializationCuboids = append(initializationCuboids, c)
	}

	normalized := initializeReactor(initializationCuboids, l)

	ct := countCubesOn(normalized)

	return strconv.FormatUint(uint64(ct), 10)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	cuboids := parseInput(r)

	l.Printf("Parsed %d cuboids from input", len(cuboids))

	normalized := initializeReactor(cuboids, l)

	l.Printf("Normalized into %d cuboids", len(normalized))

	ct := countCubesOn(normalized)

	return strconv.FormatUint(uint64(ct), 10)
}

////////////////////////////////////////////////////////////////////////////////
// Non-brute force solution

// takes a set of cuboids and returns a normalized set of non-overlapping
// cuboids that have been turned on
func initializeReactor(cuboids []cuboid, l *log.Logger) []cuboid {

	xIntervals := buildIntervals(
		cuboids,
		func(c cuboid) int {
			return c.position.x
		},
		func(c cuboid) int {
			return c.size.x
		},
	)

	yIntervals := buildIntervals(
		cuboids,
		func(c cuboid) int {
			return c.position.y
		},
		func(c cuboid) int {
			return c.size.y
		},
	)

	zIntervals := buildIntervals(
		cuboids,
		func(c cuboid) int {
			return c.position.z
		},
		func(c cuboid) int {
			return c.size.z
		},
	)

	log.Printf("%d x intervals, %d y intervals, %d z intervals (=%d combos)", len(xIntervals), len(yIntervals), len(zIntervals), len(xIntervals)*len(yIntervals)*len(zIntervals))

	result := make([]cuboid, 0)

	processed := 0
	start := time.Now()

	for _, xInterval := range xIntervals {
		for _, yInterval := range yIntervals {
			for _, zInterval := range zIntervals {

				processed++
				if processed%1000000 == 0 {
					elapsed := time.Now().Sub(start)
					msPer := float64(elapsed.Milliseconds()) / float64(processed)
					remaining := (len(xIntervals) * len(yIntervals) * len(zIntervals)) - processed
					msRemaining := float64(remaining) * msPer
					l.Printf("Processed %d intervals (%fms per; %d hit(s); ~%ds remain)", processed, msPer, len(result), int(msRemaining/1000))
				}

				cuboidIndices := intersection(xInterval.cuboidIndices, yInterval.cuboidIndices, zInterval.cuboidIndices)

				if len(cuboidIndices) == 0 {
					continue
				}

				// later `on` values in `cuboids` override earlier ones
				lastCuboidIndex := cuboidIndices[len(cuboidIndices)-1]
				on := cuboids[lastCuboidIndex].on
				if !on {
					continue
				}

				// ok so now we have x, y, and z values
				c := cuboid{
					position: point{
						x: xInterval.start,
						y: yInterval.start,
						z: zInterval.start,
					},
					size: point{
						x: xInterval.end - xInterval.start,
						y: yInterval.end - yInterval.start,
						z: zInterval.end - zInterval.start,
					},
					on: true,
				}

				result = append(result, c)
			}
		}
	}

	return result
}

func countCubesOn(cuboids []cuboid) uint {
	var ct uint

	for _, c := range cuboids {
		ct += (uint(c.size.x) * uint(c.size.y) * uint(c.size.z))
	}
	return ct
}

// given a set of cuboids along with functions to pull off specific coordinate
// + extent values (e.g. x and width, y and height, etc.), returns the set of
// intervals generated.
func buildIntervals(cuboids []cuboid, getCoordinate func(c cuboid) int, getExtent func(c cuboid) int) []interval {
	// these maps connect values to the indices of <cuboids> that contain those values
	values, valuesToCuboidIndices := buildValueMap(cuboids, func(c cuboid) []int {
		return []int{
			getCoordinate(c),
			getCoordinate(c) + getExtent(c),
		}
	})

	result := make([]interval, 0)

	// Tracks the cuboid indices that are currently "open"--we have seen their
	// `start` value but not their `end`
	openCuboidIndices := make(map[int]bool)

	var currentInterval *interval

	for _, value := range values {
		if currentInterval == nil {
			currentInterval = &interval{
				start: value,
			}

		} else {

			currentInterval.end = value
			for i, isOpen := range openCuboidIndices {
				if isOpen {
					currentInterval.cuboidIndices = append(currentInterval.cuboidIndices, i)
				}
			}
			sort.Ints(currentInterval.cuboidIndices)
			result = append(result, *currentInterval)

			currentInterval = &interval{
				start: value,
			}
		}

		for _, cuboidIndex := range valuesToCuboidIndices[value] {
			openCuboidIndices[cuboidIndex] = !openCuboidIndices[cuboidIndex]
		}
	}

	return result
}

// given a set of cuboids and a function to read coordinate values them, returns
// a sorted slice of those values and a map connecting those values to the
// cuboid indices that contained those values
func buildValueMap(cuboids []cuboid, f func(c cuboid) []int) ([]int, map[int][]int) {
	valueToIndexMap := make(map[int][]int)

	for cuboidIndex, c := range cuboids {
		for _, value := range f(c) {
			index := sort.SearchInts(valueToIndexMap[value], cuboidIndex)
			alreadyPresent := index < len(valueToIndexMap[value]) && valueToIndexMap[value][index] == cuboidIndex
			if !alreadyPresent {
				valueToIndexMap[value] = append(valueToIndexMap[value], cuboidIndex)
				// TODO: Actually insert at the right place
				sort.Ints(valueToIndexMap[value])
			}
		}
	}

	values := make([]int, 0, len(valueToIndexMap))
	for value := range valueToIndexMap {
		values = append(values, value)
	}
	sort.Ints(values)

	return values, valueToIndexMap
}

func intersection(a, b, c []int) []int {

	// these slices are sorted
	// if the min of any slice > max of any slice, there is no intersection
	// if the max of any slice < min of any slice, there is no intersection

	aLen := len(a)
	bLen := len(b)
	cLen := len(c)

	aMin := a[0]
	aMax := a[aLen-1]
	bMin := b[0]
	bMax := b[bLen-1]
	cMin := c[0]
	cMax := c[cLen-1]

	minLen := aLen
	if bLen < minLen {
		minLen = bLen
	}
	if cLen < minLen {
		minLen = cLen
	}

	if aMin > bMax || aMin > cMax {
		return []int{}
	}

	if bMin > aMax || bMin > cMax {
		return []int{}
	}

	if cMin > aMax || cMin > bMax {
		return []int{}
	}

	result := make([]int, 0, minLen)

	for _, aValue := range a {
		if aValue > bMax || aValue > cMax {
			break
		}

		for _, bValue := range b {

			if bValue > aMax || bValue > cMax {
				break
			}

			if bValue > aValue {
				break
			}

			if bValue != aValue {
				continue
			}

			for _, cValue := range c {
				if cValue > aMax || cValue > bMax {
					break
				}

				if cValue > aValue {
					break
				}

				if cValue == aValue {
					result = append(result, aValue)
				}
			}

		}

	}

	return result

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
