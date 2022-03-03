package d25

import (
	"log"
	"strings"
	"testing"
)

func TestExample(t *testing.T) {
	input := strings.TrimSpace(`
v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>
	`)
	solution := Puzzle1(strings.NewReader(input), log.Default())
	if solution != "58" {
		t.Errorf("Expected %d, got %s", 58, solution)
	}
}
