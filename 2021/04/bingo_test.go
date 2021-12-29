package day04

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewGame(t *testing.T) {
	input := `
1,2,3,4

1 2
3 4

5 6
7 8
`

	game := NewGame(strings.NewReader(input))

	assertIntSlicesEqual(t, []int{1, 2, 3, 4}, game.numbers)

	assertIntSlicesEqual(t, []int{1, 2}, getSquareValues(game.boards[0].squares[0]))

}

func TestRunGame(t *testing.T) {
	input := `
7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
8  2 23  4 24
21  9 14 16  7
6 10  3 18  5
1 12 20 15 19

3 15  0  2 22
9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
2  0 12  3  7
	`

	game := NewGame(strings.NewReader(input))

	assertIntSlicesEqual(
		t,
		game.numbers,
		[]int{7, 4, 9, 5, 11, 17, 23, 2, 0, 14, 21, 24, 10, 16, 13, 6, 15, 25, 12, 22, 18, 20, 8, 19, 3, 26, 1},
	)

	assertIntSlicesEqual(
		t,
		getSquareValues(game.boards[0].squares[0]),
		[]int{22, 13, 17, 11, 0},
	)

	solution := game.Run()

	firstWinner := solution[0]
	if firstWinner.index != 3 {
		fmt.Println(firstWinner.String())
		t.Error("The wrong board won first")
	}

	lastWinner := solution[len(solution)-1]
	if lastWinner.index != 2 {
		fmt.Println(lastWinner.String())
		t.Error("the wrong board was last to win")
	}

	if lastWinner.finalScore != 1924 {
		fmt.Println(lastWinner.String())
		t.Error("the last board had the wrong score")
	}

}

func assertIntSlicesEqual(t *testing.T, x []int, y []int) {
	if len(x) != len(y) {
		t.Error(fmt.Sprintf("Lengths don't match: %d vs %d", len(x), len(y)))
		return
	}

	for i := range x {
		if x[i] != y[i] {
			t.Error(fmt.Sprintf("Values don't match at position %d: %d != %d", i, x[i], y[i]))
		}
	}
}

func getSquareValues(squares []square) []int {
	var result []int
	for i := range squares {
		result = append(result, squares[i].value)
	}
	return result
}
