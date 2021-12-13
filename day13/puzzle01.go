package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type foldInstruction struct {
	x, y int
}

type sheet struct {
	dots         []point
	instructions []foldInstruction
}

func main() {

	sheet := parseInput(os.Stdin)

	for len(sheet.instructions) > 0 {
		i := sheet.instructions[0]
		nextSheet := sheet.fold()

		axis := "x"
		value := i.x
		if i.y > 0 {
			axis = "y"
			value = i.y
		}

		fmt.Printf("fold along %s=%d leaves %d points\n", axis, value, len(nextSheet.dots))
		sheet = nextSheet
	}

	fmt.Println("Final code:")
	fmt.Println(sheet.String())

}

// fold processes the next fold instruction and returns a new sheet
func (s *sheet) fold() sheet {
	if len(s.instructions) == 0 {
		return *s
	}

	instruction := s.instructions[0]
	result := sheet{
		instructions: s.instructions[1:],
	}

	for _, p := range s.dots {

		foldedP := p

		if instruction.x > 0 {
			// fold is along the x axis, result will be narrower
			if p.x > instruction.x {
				foldedP.x = instruction.x - (p.x - instruction.x)
			}
		} else if instruction.y > 0 {
			// fold is along the y axis, result will be shorter
			if p.y > instruction.y {
				foldedP.y = instruction.y - (p.y - instruction.y)
			}
		}

		// now we make sure our resulting dots are unique
		alreadyThere := false
		for _, o := range result.dots {
			if o == foldedP {
				alreadyThere = true
				break
			}
		}
		if !alreadyThere {
			result.dots = append(result.dots, foldedP)
		}
	}

	return result
}

func (s *sheet) String() string {

	// step 1 = make the final set of dots
	var maxX, maxY int
	for _, d := range s.dots {
		if d.x > maxX {
			maxX = d.x
		}
		if d.y > maxY {
			maxY = d.y
		}
	}

	width := maxX + 1
	height := maxY + 1

	grid := make([][]rune, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			grid[y][x] = ' '
		}
	}

	for _, d := range s.dots {
		grid[d.y][d.x] = 'X'
	}

	b := strings.Builder{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b.WriteRune(grid[y][x])
		}
		if y < height-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func parseInput(r io.Reader) sheet {
	s := bufio.NewScanner(r)
	foldRx := regexp.MustCompile("fold along (x|y)=(\\d+)")

	result := sheet{}

	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if len(l) == 0 {
			continue
		}

		parts := strings.Split(l, ",")
		if len(parts) == 2 {
			result.dots = append(result.dots, parsePoint(parts[0], parts[1]))
			continue
		}

		m := foldRx.FindStringSubmatch(l)
		if m == nil {
			continue
		}

		pos, err := strconv.ParseInt(m[2], 10, 32)
		if err != nil {
			panic(err)
		}

		if m[1] == "x" {
			result.instructions = append(result.instructions, foldInstruction{int(pos), 0})
		} else if m[1] == "y" {
			result.instructions = append(result.instructions, foldInstruction{0, int(pos)})
		}
	}

	return result
}

func parsePoint(xStr, yStr string) point {
	x, xErr := strconv.ParseInt(xStr, 10, 32)
	if xErr != nil {
		panic("Invalid X")
	}
	y, yErr := strconv.ParseInt(yStr, 10, 32)
	if yErr != nil {
		panic("Invalid Y")
	}
	return point{int(x), int(y)}
}
