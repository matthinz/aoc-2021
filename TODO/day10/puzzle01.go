package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var ChunkDelimiters = map[rune]rune{
	'(': ')',
	'[': ']',
	'{': '}',
	'<': '>',
}

var IllegalCharacterScores = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var CompletionCharacterScores = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

type LineState int

const (
	LineStateValid LineState = iota
	LineStateIncomplete
	LineStateCorrupted
)

type parsedLine struct {
	text       string
	state      LineState
	errorRune  rune
	errorPos   int
	completion string
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	var lineNumber int
	var errorScore int

	var incompleteLines []parsedLine

	for s.Scan() {
		lineNumber++
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		p := readLine(line)

		if p.state == LineStateCorrupted {
			errorScore += IllegalCharacterScores[p.errorRune]
		} else if p.state == LineStateIncomplete {
			incompleteLines = append(incompleteLines, p)
		}
	}

	fmt.Printf("Total error score: %d\n", errorScore)

	sort.Slice(incompleteLines, func(i, j int) bool {
		return incompleteLines[i].completionScore() < incompleteLines[j].completionScore()
	})

	// 7
	// 7 / 2 = 3

	midIndex := len(incompleteLines) / 2

	for index, i := range incompleteLines {
		var arrow string
		if index == midIndex {
			arrow = " <---"
		}
		fmt.Printf("%s : %d%s\n", i.completion, i.completionScore(), arrow)
	}

	fmt.Printf("Middle completion score: %d\n", incompleteLines[midIndex].completionScore())

}

func (p *parsedLine) completionScore() int {
	var result int

	for _, r := range p.completion {
		result *= 5
		result += CompletionCharacterScores[r]
	}

	return result
}

func (s LineState) String() string {
	switch s {
	case LineStateValid:
		return "valid"
	case LineStateCorrupted:
		return "corrupted"
	case LineStateIncomplete:
		return "incomplete"
	default:
		return "unknown"
	}
}

func readLine(line string) parsedLine {

	var stack []rune

	for pos, r := range line {
		_, found := ChunkDelimiters[r]

		if found {
			// <r> was opening a new chunk
			stack = append(stack, r)
			continue
		}

		// <r> is _hopefully_ closing a chunk
		if len(stack) == 0 {
			panic("nothing in the stack")
		}

		lastChunk := stack[len(stack)-1]
		expectedDelimiter := ChunkDelimiters[lastChunk]

		if r != expectedDelimiter {
			// this is is corrupted
			return parsedLine{
				text:      line,
				state:     LineStateCorrupted,
				errorRune: r,
				errorPos:  pos,
			}
		}

		// TODO: better way to grab all but the last thing?
		stack = stack[0 : len(stack)-1]
	}

	if len(stack) > 0 {
		return parsedLine{
			text:  line,
			state: LineStateIncomplete,
			completion: mapString(
				reverseString(string(stack)),
				func(r rune) rune {
					return ChunkDelimiters[r]
				},
			),
		}
	}

	return parsedLine{
		text:  line,
		state: LineStateValid,
	}
}

func mapString(input string, f func(r rune) rune) string {
	result := make([]rune, 0, len(input))
	for _, r := range input {
		result = append(result, f(r))
	}
	return string(result)
}

func reverseString(input string) string {
	inputLen := len(input)
	result := make([]rune, inputLen)
	for i, r := range input {
		result[inputLen-1-i] = r
	}
	return string(result)
}
