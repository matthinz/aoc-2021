package d13

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func parseInput(r io.Reader) sheet {

	foldRx := regexp.MustCompile("fold along (x|y)=(\\d+)")
	s := bufio.NewScanner(r)

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
