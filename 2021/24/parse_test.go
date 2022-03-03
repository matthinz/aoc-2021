package d24

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
)

//go:embed input
var realInput string

func TestParseFirstLinesOfRealInput(t *testing.T) {
	t.Skip()
	const LineCount = 170
	lines := strings.Split(realInput, "\n")
	first := strings.Join(lines[0:LineCount], "\n")
	parseInput(strings.NewReader(first))
}

func TestParseRealInputFindsAllInputsInZ(t *testing.T) {
	r := parseInput(strings.NewReader(realInput))
	inputsFound := make(map[int]int)
	r.z.Accept(func(e Expression) {
		ie, ok := e.(*InputExpression)
		if ok {
			inputsFound[ie.index]++
		}
	})

	for i := 0; i < 14; i++ {
		_, ok := inputsFound[i]
		if !ok {
			t.Errorf("Input not found: %d", i)
		}
	}
}

func TestParseRealInputTrySolution(t *testing.T) {
	t.Skip()
	registers := parseInput(strings.NewReader(realInput))
	solution := []int{9, 8, 7, 1, 4, 3, 9, 3, 4, 9, 7, 9, 3, 3}

	simplified := registers.z.Simplify(solution)
	r := simplified.Range()

	if !r.Includes(0) {
		t.Fatalf("Solution range does not include 0")
	}

	fmt.Println(simplified.Range())

	value, err := simplified.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if value != 0 {
		t.Errorf("solution evaluated to %d", value)
	}

}
