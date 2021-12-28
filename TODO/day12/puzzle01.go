package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type cave struct {
	name string
	big  bool
}

type connection [2]cave

func main() {
	connections := parseInput(os.Stdin)

	for _, c := range connections {
		fmt.Printf("%s <-> %s\n", c[0].String(), c[1].String())
	}

	initialPath := []cave{
		cave{
			name: "start",
			big:  false,
		},
	}

	paths := buildPaths(connections, initialPath, 1)

	fmt.Printf("%d paths found visiting small caves only once\n", len(paths))

	paths = buildPaths(connections, initialPath, 2)

	fmt.Printf("%d paths found when we're allowed to revisit one small cave once\n", len(paths))

}

func NewCave(name string) cave {
	return cave{
		name: name,
		big:  strings.ToUpper(name) == name,
	}
}

func (c *cave) String() string {
	if c.big {
		return fmt.Sprintf("%s (big)", c.name)
	} else {
		return fmt.Sprintf("%s (small)", c.name)
	}
}

func (c *cave) VisitCount(path []cave) int {
	var result int
	for _, x := range path {
		if (*c) == x {
			result++
		}
	}
	return result
}

func buildPaths(connections []connection, currentPath []cave, puzzle int) [][]cave {

	lastCave := currentPath[len(currentPath)-1].name

	nextSteps := make([]cave, 0)
	for _, c := range connections {
		if c[0].name == lastCave && canVisit(c[1], currentPath, puzzle) {
			nextSteps = append(nextSteps, c[1])
		} else if c[1].name == lastCave && canVisit(c[0], currentPath, puzzle) {
			nextSteps = append(nextSteps, c[0])
		}
	}

	var result [][]cave

	if len(nextSteps) == 0 {
		return result
	}

	for _, s := range nextSteps {

		newPath := make([]cave, len(currentPath)+1)
		copy(newPath, currentPath)
		newPath[len(newPath)-1] = s

		if s.name == "end" {
			// we have a completed path
			result = append(result, newPath)
		} else {
			for _, p := range buildPaths(connections, newPath, puzzle) {
				result = append(result, p)
			}
		}
	}

	return result
}

func canVisit(c cave, path []cave, puzzle int) bool {
	if c.big {
		// we can *always* visit big caves more than once
		return true
	}

	visitCount := c.VisitCount(path)

	if puzzle == 1 {
		// small caves get 1 visit
		return visitCount < 1
	}

	// For puzzle 2, the rules are more complicated:

	if visitCount < 1 {
		return true
	}

	if (c.name == "start" || c.name == "end") && visitCount > 0 {
		// can't ever revisit start or end
		return false
	}

	// we are allowed to grant 1 small cave 2 visits
	grantedExtraVisit := false

	for i := 0; i < len(path); i++ {
		if path[i].big {
			continue
		}
		if path[i].VisitCount(path) > 1 {
			grantedExtraVisit = true
			break
		}
	}

	return !grantedExtraVisit
}

func parseInput(r io.Reader) []connection {
	b := bufio.NewScanner(r)

	result := make([]connection, 0)

	for b.Scan() {
		line := strings.TrimSpace(b.Text())
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}

		result = append(result, connection{NewCave(parts[0]), NewCave(parts[1])})
	}

	return result
}
