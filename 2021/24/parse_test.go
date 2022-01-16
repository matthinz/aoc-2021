package d24

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
)

//go:embed input
var realInput string

func TestParseRealInputFindsAllInputsInZ(t *testing.T) {
	r := parseInput(strings.NewReader(realInput))
	inputsFound := make(map[int]int)
	r.z.Accept(func(e Expression) {
		ie, ok := e.(*InputExpression)
		if ok {
			fmt.Println(ie)
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
