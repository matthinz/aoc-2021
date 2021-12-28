package d13

import "strings"

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
