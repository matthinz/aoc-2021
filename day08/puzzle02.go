package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type input struct {
	signalValues []string
	outputValues []string
}

var SignalValuesForNumbers = [10]string{
	"abcefg",
	"cf",
	"acdeg",
	"acdfg",
	"bcdf",
	"abdfg",
	"abdefg",
	"acf",
	"abcdefg",
	"abcdfg",
}

func main() {
	inputs := parseInput(os.Stdin)
	ch := make(chan int)

	for _, i := range inputs {
		go func(i input) {
			key := interpretSignalValues(i.signalValues)
			digits := getDigits(i, key)
			ch <- makeIntFromDigits(digits)
		}(i)
	}

	var received int
	var sum int

	for num := range ch {
		received++

		sum += num
		fmt.Printf("%d (total %d) - (%d / %d)\n", num, sum, received, len(inputs))

		if received == len(inputs) {
			close(ch)
			break
		}
	}

}

func makeIntFromDigits(digits []int) int {
	var result int
	for i, digit := range digits {
		result += (digit * int(math.Pow10(len(digits)-(i+1))))
	}
	return result
}

func getDigits(i input, key []int) []int {

	result := make([]int, len(i.outputValues))

	for outputIndex, v := range i.outputValues {
		inputIndex := stringSliceIndexOf(i.signalValues, v)
		if inputIndex < 0 {
			panic(fmt.Sprintf("output %s not found in signal values %v", v, i.signalValues))
		}
		digit := key[inputIndex]
		result[outputIndex] = digit
	}

	return result

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
			signalValues: parseSegments(parts[0]),
			outputValues: parseSegments(parts[1]),
		})
	}

	return result
}

func parseSegments(input string) []string {
	input = strings.Trim(input, " \t\r\n")
	segments := strings.Split(input, " ")

	result := make([]string, len(segments))
	for i, segment := range segments {
		runes := []rune(segment)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		result[i] = string(runes)
	}
	return result
}

func interpretSignalValues(inputs []string) []int {

	to := "abcdefg"
	var solution string
	var ct int

	for from := range generateStrings(len(to)) {
		ct++
		ok := areMappingsValid(from, to, inputs)
		if ok {
			solution = from
			break
		}
	}

	if len(solution) == 0 {
		return []int{}
	}

	// result is a slice of ints where each number corresponds to the input
	result := make([]int, len(inputs))

	for i, input := range inputs {
		transformed := applyTransform(input, solution, to)
		num := stringSliceIndexOf(SignalValuesForNumbers[:], transformed)
		if num < 0 {
			panic(fmt.Sprintf("solution %s does not solve %s", solution, input))
		}
		result[i] = num
	}

	return result
}

func buildStrings(targetLength int, alphabet string, ch chan string, input string) {
	if len(input) == targetLength {
		ch <- input
		return
	}
	for _, c := range alphabet {
		buildStrings(targetLength, strings.Replace(alphabet, string(c), "", 1), ch, input+string(c))
	}
}

func generateStrings(targetLength int) chan string {
	ch := make(chan string)

	go func() {
		buildStrings(targetLength, "abcdefg", ch, "")
		close(ch)
	}()

	return ch
}

// given an input string, converts each occurrence of a character in
// <from> to the corresponding character in <to>
func applyTransform(input string, from string, to string) string {

	result := make([]rune, len(input))

	for i, r := range input {
		fromIndex := strings.IndexRune(from, r)
		if fromIndex >= 0 {
			result[i] = rune(to[fromIndex])
		} else {
			result[i] = r
		}
	}

	// sort the result a->z
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })

	return string(result)
}

func areMappingsValid(from string, to string, inputs []string) bool {

	// in each input, change each character in <from> to the corresponding
	// character in <to>

	transformedInputs := make([]string, len(inputs))
	for i, input := range inputs {
		transformedInputs[i] = applyTransform(input, from, to)
	}

	// having been transformed, verify that each new input is present
	// in our known good set of numbers

	for _, transformedInputSignalValue := range transformedInputs {
		found := false
		for _, knownGoodSignalValue := range SignalValuesForNumbers {
			if transformedInputSignalValue == knownGoodSignalValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func stringSliceIndexOf(slice []string, value string) int {
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

func stringSliceContains(slice []string, value string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}

func replaceRunesAndSort(input string, from []rune, to []rune) string {
	if len(from) != len(to) {
		panic(fmt.Sprintf("wrong length (%d vs %d)", len(from), len(to)))
	}
	result := make([]rune, len(input))
	for i, c := range input {
		for j := range from {
			if c == from[j] {
				result[i] = to[j]
			} else {
				result[i] = c
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })

	return string(result)
}

func findStringsWithLength(values []string, length int) []string {
	var result []string
	for _, value := range values {
		if len(value) == length {
			result = append(result, value)
		}
	}
	return result
}

func findValuesWithSegment(signalValues []string, segment rune) []string {
	resultMap := make(map[string]bool)
	for _, value := range signalValues {
		if strings.ContainsRune(value, segment) {
			resultMap[value] = true
		}
	}

	result := make([]string, 0, len(resultMap))
	for value := range resultMap {
		result = append(result, value)
	}
	return result
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func remove(values []string, index int) []string {
	result := make([]string, 0, len(values)-1)
	for i, value := range values {
		if i != index {
			result = append(result, value)
		}
	}
	return result
}
