package d24

import (
	"fmt"
	"strconv"
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

func TestFirstExample(t *testing.T) {

	input := strings.TrimSpace(`
inp x
mul x -1`)
	reg := parseInput(strings.NewReader(input))

	tests := map[int]int{
		0:    0,
		1:    -1,
		-1:   1,
		100:  -100,
		-100: 100,
	}

	for input, expected := range tests {
		t.Run(strconv.Itoa(input), func(t *testing.T) {
			inputs := []int{input}
			expr := reg.x.Simplify(inputs)
			actual, err := expr.Evaluate()
			if err != nil {
				t.Fatal(err)
			}
			if actual != expected {
				t.Errorf("Expected %d, but got %d", expected, actual)
			}
		})
	}
}

func TestSecondExample(t *testing.T) {

	input := strings.TrimSpace(`
inp z
inp x
mul z 3
eql z x`)
	reg := parseInput(strings.NewReader(input))

	tests := map[[2]int]int{
		{1, 3}: 1,
		{1, 4}: 0,
		{2, 6}: 1,
		{2, 9}: 0,
	}

	for inputs, expected := range tests {
		t.Run(fmt.Sprintf("%v", inputs), func(t *testing.T) {
			expr := reg.z.Simplify(inputs[:])
			actual, err := expr.Evaluate()
			if err != nil {
				t.Fatal(err)
			}
			if actual != expected {
				t.Errorf("Expected %d, but got %d", expected, actual)
			}
		})
	}
}
