package d18

import (
	"fmt"
	"testing"
)

func TestExplode(t *testing.T) {

	input := "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]"
	number, err := parseSnailfishNumber(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	p := number.left // [3,[2,[1,[7,3]]]]
	if p.String() != "[3,[2,[1,[7,3]]]]" {
		t.Fatalf("number.left was wrong")
	}

	p = p.right // [2,[1,[7,3]]]
	if p.String() != "[2,[1,[7,3]]]" {
		t.Fatalf("number.right was wrong")
	}

	p = p.right // [1,[7,3]]
	if p.String() != "[1,[7,3]]" {
		t.Fatalf("number.right was wrong")
	}

	p = p.right // [7,3]
	if p.String() != "[7,3]" {
		t.Fatalf("number.right was wrong")
	}

	p.explode()

	expected := "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]"
	actual := number.String()

	if actual != expected {
		t.Fatalf("explode() failed. expected %s, got %s", expected, actual)
	}
}

func TestReduce(t *testing.T) {
	input := "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]"
	number, err := parseSnailfishNumber(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	number.reduce()

	expected := "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"
	actual := number.String()

	if actual != expected {
		t.Fatalf("reduce() failed. expected %s, got %s", expected, actual)
	}
}

func TestSplit(t *testing.T) {

	input := "[5,8]"
	number, err := parseSnailfishNumber(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	n := number.left
	n.split()
	expected := "[2,3]"
	actual := n.String()
	if actual != expected {
		t.Fatalf("split() failed. expected %s, got %s", expected, actual)
	}

	n = number.right
	n.split()
	expected = "[4,4]"
	actual = n.String()
	if actual != expected {
		t.Fatalf("split() failed. expected %s, got %s", expected, actual)
	}

	expected = "[[2,3],[4,4]]"
	actual = number.String()
	if actual != expected {
		t.Fatalf("split() failed. expected %s, got %s", expected, actual)
	}
}

func TestSum(t *testing.T) {
	input := []string{
		"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
		"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
		"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
		"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
		"[7,[5,[[3,8],[1,4]]]]",
		"[[2,[2,2]],[8,[8,1]]]",
		"[2,9]",
		"[1,[[[9,3],9],[[9,0],[0,7]]]]",
		"[[[5,[7,4]],7],1]",
		"[[[[4,2],2],6],[8,7]]",
	}

	expected := []string{
		"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
		"[[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]",
		"[[[[6,7],[6,7]],[[7,7],[0,7]]],[[[8,7],[7,7]],[[8,8],[8,0]]]]",
		"[[[[7,0],[7,7]],[[7,7],[7,8]]],[[[7,7],[8,8]],[[7,7],[8,7]]]]",
		"[[[[7,7],[7,8]],[[9,5],[8,7]]],[[[6,8],[0,8]],[[9,9],[9,0]]]]",
		"[[[[6,6],[6,6]],[[6,0],[6,7]]],[[[7,7],[8,9]],[8,[8,1]]]]",
		"[[[[6,6],[7,7]],[[0,7],[7,7]]],[[[5,5],[5,6]],9]]",
		"[[[[7,8],[6,7]],[[6,8],[0,8]]],[[[7,7],[5,0]],[[5,5],[5,6]]]]",
		"[[[[7,7],[7,7]],[[8,7],[8,7]]],[[[7,0],[7,7]],9]]",
		"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
	}

	var num *snailfishNumber

	for i, s := range input {
		n, err := parseSnailfishNumber(s)
		if err != nil {
			t.Fatal(err.Error())
		}

		if num == nil {
			num = n
		} else {
			fmt.Printf("%s + %s\n", num.String(), n.String())
			num = num.add(*n)
		}

		num.checkDepths(0)

		if num.String() != expected[i] {
			t.Fatalf("Step %d failed: expected %s, got %s", i, expected[i], num.String())
		}

	}

}
