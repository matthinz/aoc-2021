package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x int
	y int
}

type line struct {
	start point
	end   point
}

func main() {

	lines := ParseInput(os.Stdin)

	intersections := CalculateIntersections(lines, true)

	var atLeast2 int

	for _, row := range intersections {
		for _, ct := range row {
			if ct >= 2 {
				atLeast2++
			}
		}
	}

	fmt.Printf("Points w/ at least 2 overlaps (including diagonals): %d\n", atLeast2)

}

func (l *line) containsPoint(p point) bool {

	var minX, minY, maxX, maxY int

	if l.start.x < l.end.x {
		minX = l.start.x
		maxX = l.end.x
	} else {
		minX = l.end.x
		maxX = l.start.x
	}

	if l.start.y < l.end.y {
		minY = l.start.y
		maxY = l.end.y
	} else {
		minY = l.end.y
		maxY = l.start.y
	}

	xInBounds := p.x >= minX && p.x <= maxX
	yInBounds := p.y >= minY && p.y <= maxY

	if !(xInBounds && yInBounds) {
		return false
	}

	var onTheLine bool

	if l.isHorizontal() {
		onTheLine = p.y == l.start.y
	} else if l.isVertical() {
		onTheLine = p.x == l.start.x
	} else {
		slope := (l.end.y - l.start.y) / (l.end.x - l.start.x)
		yIntercept := l.start.y - (slope * l.start.x)
		onTheLine = (p.y == slope*p.x+yIntercept)
	}

	if p.x == 1 && p.y == 9 && onTheLine {
		fmt.Println(l)
	}

	return onTheLine
}

func (l *line) isHorizontal() bool {
	return l.start.y == l.end.y
}

func (l *line) isVertical() bool {
	return l.start.x == l.end.x
}

func CalculateIntersections(lines []line, includeDiagonals bool) [][]int {
	var candidateLines []line

	if includeDiagonals {
		candidateLines = lines
	} else {
		for i := range lines {
			if lines[i].isHorizontal() || lines[i].isVertical() {
				candidateLines = append(candidateLines, lines[i])
			}
		}
	}

	min, max := getMinMaxPoints(candidateLines)

	width := max.x - min.x + 1
	height := max.y - min.y + 1

	result := make([][]int, height)

	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			p := point{x + min.x, y + min.y}
			for i := range candidateLines {
				if candidateLines[i].containsPoint(p) {
					result[y][x]++
				}
			}
		}
	}

	return result
}

func getMinMaxPoints(lines []line) (point, point) {
	var min, max point
	for i := range lines {
		for _, p := range []point{lines[i].start, lines[i].end} {
			if p.x < min.x {
				min.x = p.x
			}
			if p.y < min.y {
				min.y = p.y
			}
			if p.x > max.x {
				max.x = p.x
			}
			if p.y > max.y {
				max.y = p.y
			}
		}
	}
	return min, max
}

func ParseInput(r io.Reader) []line {
	scanner := bufio.NewScanner(r)

	var result []line

	for scanner.Scan() {
		line := parseLine(scanner.Text())
		if line == nil {
			continue
		}

		result = append(result, *line)
	}

	return result
}

func parseLine(input string) *line {

	rx, err := regexp.Compile("(\\d+),(\\d+)\\s*->\\s*(\\d+),(\\d+)")
	if err != nil {
		panic(err)
	}

	m := rx.FindStringSubmatch(input)

	if m == nil {
		return nil
	}

	nums := [4]int{
		0, 0, 0, 0,
	}

	for i := range nums {
		token := m[i+1]
		value, err := strconv.ParseInt(token, 10, 32)
		if err != nil {
			panic(err)
		}
		nums[i] = int(value)
	}

	return &line{
		start: point{nums[0], nums[1]},
		end:   point{nums[2], nums[3]},
	}
}
