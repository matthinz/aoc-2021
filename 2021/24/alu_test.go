package d24

import (
	"strings"
	"testing"
)

func TestFindAllInputsInZ(t *testing.T) {
	t.Skip()
	reg := parseInput(strings.NewReader(realInput))
	expr := reg.z.Simplify([]int{})
	inputs := make(map[int]int)
	expr.Accept(func(e Expression) {
		if input, isInput := e.(*InputExpression); isInput {
			inputs[input.index]++
		}
	})
	if len(inputs) != 14 {
		t.Errorf("Expected %d inputs, but got %d", 14, len(inputs))
	}
	for i := 0; i < 14; i++ {
		_, ok := inputs[i]
		if !ok {
			t.Errorf("Input %d not found", i)
		}
	}
}

func BenchmarkSimplify(b *testing.B) {
	reg := parseInput(strings.NewReader(realInput))
	for i := 0; i < b.N; i++ {
		reg.z.Simplify([]int{})
	}
}
