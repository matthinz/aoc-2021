package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const INPUT = `
0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2
	`

func TestParseInput(t *testing.T) {

	expected := []line{
		line{point{0, 9}, point{5, 9}},
		line{point{8, 0}, point{0, 8}},
		line{point{9, 4}, point{3, 4}},
		line{point{2, 2}, point{2, 1}},
		line{point{7, 0}, point{7, 4}},
		line{point{6, 4}, point{2, 0}},
		line{point{0, 9}, point{2, 9}},
		line{point{3, 4}, point{1, 4}},
		line{point{0, 0}, point{8, 8}},
		line{point{5, 5}, point{8, 2}},
	}

	lines := ParseInput(strings.NewReader(INPUT))

	if len(lines) != 10 {
		t.Error(fmt.Sprintf("Expected 10 lines, got %d", len(lines)))
	}

	for i := range expected {
		if lines[i] != expected[i] {
			t.Error(fmt.Sprintf("%d: expected %v, actual %v", i, expected[i], lines[i]))
		}
	}

}

func TestContainsPoint(t *testing.T) {

	l := line{start: point{0, 9}, end: point{5, 9}}

	if !l.containsPoint(point{1, 9}) {
		t.Error("should've contained point")
	}

	if !l.containsPoint(point{5, 9}) {
		t.Error("should've contained point")
	}

	if l.containsPoint(point{6, 9}) {
		t.Error("should not contain point")
	}

	if l.containsPoint(point{1, 1}) {
		t.Error("should not contain point")
	}

}

func TestCalculateIntersectionsNoDiagonals(t *testing.T) {

	lines := ParseInput(strings.NewReader(INPUT))

	if len(lines) != 10 {
		t.Error("wrong # of lines parsed")
		return
	}

	intersections := CalculateIntersections(lines, false)

	expectedIntersections := [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		[]int{0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
		[]int{0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		[]int{0, 1, 1, 2, 1, 1, 1, 2, 1, 1},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{2, 2, 2, 1, 1, 1, 0, 0, 0, 0},
	}

	ok := len(intersections) == len(expectedIntersections)

	if ok {
		for y := range expectedIntersections {
			expected := expectedIntersections[y]
			actual := intersections[y]
			if len(expected) != len(actual) {
				ok = false
				continue
			}
			for x := range expected {
				if expected[x] != actual[x] {
					ok = false
				}
			}
		}
	}

	if !ok {
		fmt.Println("EXPECTED")
		fmt.Println(formatIntersections(expectedIntersections))
		fmt.Println("ACTUAL")
		fmt.Println(formatIntersections(intersections))
		t.Error("Intersections are wrong")
	}

}

func TestCalculateIntersectionsWithDiagonals(t *testing.T) {

	lines := ParseInput(strings.NewReader(INPUT))

	if len(lines) != 10 {
		t.Error("wrong # of lines parsed")
		return
	}

	intersections := CalculateIntersections(lines, true)

	expectedIntersections := [][]int{
		[]int{1, 0, 1, 0, 0, 0, 0, 1, 1, 0},
		[]int{0, 1, 1, 1, 0, 0, 0, 2, 0, 0},
		[]int{0, 0, 2, 0, 1, 0, 1, 1, 1, 0},
		[]int{0, 0, 0, 1, 0, 2, 0, 2, 0, 0},
		[]int{0, 1, 1, 2, 3, 1, 3, 2, 1, 1},
		[]int{0, 0, 0, 1, 0, 2, 0, 0, 0, 0},
		[]int{0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		[]int{0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
		[]int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		[]int{2, 2, 2, 1, 1, 1, 0, 0, 0, 0},
	}

	ok := len(intersections) == len(expectedIntersections)

	if ok {
		for y := range expectedIntersections {
			expected := expectedIntersections[y]
			actual := intersections[y]
			if len(expected) != len(actual) {
				ok = false
				continue
			}
			for x := range expected {
				if expected[x] != actual[x] {
					ok = false
				}
			}
		}
	}

	if !ok {
		fmt.Println("EXPECTED")
		fmt.Println(formatIntersections(expectedIntersections))
		fmt.Println("ACTUAL")
		fmt.Println(formatIntersections(intersections))
		t.Error("Intersections are wrong")
	}

}

func formatIntersections(intersections [][]int) string {
	b := strings.Builder{}
	for _, row := range intersections {
		b.WriteString("\n")
		for _, value := range row {
			if value > 0 {
				b.WriteString(strconv.FormatInt(int64(value), 10))

			} else {
				b.WriteString(".")
			}
		}
	}
	return b.String()
}
