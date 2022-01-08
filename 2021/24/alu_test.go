package d24

import (
	"strings"
	"testing"
)

func TestAlu(t *testing.T) {
	input := `
inp w
add z w
mod z 2
div w 2
add y w
mod y 2
div w 2
add x w
mod x 2
div w 2
mod w 2
	`
	ops := parseInput(strings.NewReader(input))

	if len(ops) != 11 {
		t.Fatalf("Expected 11 ops, got %d", len(ops))
	}

	a := executeAll([]int{10}, ops)

	if a.w != 1 {
		t.Fatalf("expected w to be 1, but was %d", a.w)
	}

	if a.x != 0 {
		t.Fatalf("expected x to be 0, but was %d", a.x)
	}

	if a.y != 1 {
		t.Fatalf("expected y to be 1, but was %d", a.y)
	}

	if a.z != 0 {
		t.Fatalf("expected z to be 1, but was %d", a.z)
	}

}
