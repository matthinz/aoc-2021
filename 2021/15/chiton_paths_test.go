package d15

import (
	"log"
	"strings"
	"testing"
)

func TestSolveDijkstra(t *testing.T) {
	input := strings.NewReader(strings.TrimSpace(`
1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581
	`))

	grid := parseInput(input)

	lowestTotalRisk := solveDijkstra(grid, log.Default())
	if lowestTotalRisk != 40 {
		t.Fatalf("Wrong answer -- expected %d, got %d", 40, lowestTotalRisk)
	}

}
func TestInflateGrid(t *testing.T) {
	input := strings.NewReader(strings.TrimSpace(`
1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581
	`))

	grid := parseInput(input)

	inflated := inflateGrid(grid, 5)

	if (*inflated)[0][10] != 2 {
		t.Fatalf("inflate failed")
	}

	if (*inflated)[49][49] != 9 {
		t.Fatalf("inflate failed")
	}

}
