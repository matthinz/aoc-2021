package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type input struct {
	signalPatterns []string
	outputValues   []string
}

func main() {
	inputs := parseInput(os.Stdin)
	var ct int
	for _, i := range inputs {
		for _, val := range i.outputValues {
			switch len(val) {
			case 2:
				ct++
				break
			case 3:
				ct++
				break
			case 4:
				ct++
				break
			case 7:
				ct++
				break
			}
		}
	}
	fmt.Println(ct)
}

func parseInput(r io.Reader) []input {

	var result []input

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		result = append(result, input{
			signalPatterns: parseSegments(parts[0]),
			outputValues:   parseSegments(parts[1]),
		})
	}

	return result
}

func parseSegments(input string) []string {
	input = strings.Trim(input, " \t\r\n")
	return strings.Split(input, " ")
}
