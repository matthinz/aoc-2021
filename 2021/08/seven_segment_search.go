package d08

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
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

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(8, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	inputs := parseInput(r)
	ch := make(chan []int)

	for _, i := range inputs {
		go func(i input) {
			key := interpretSignalValues(i.signalValues)
			ch <- getDigits(i, key)
		}(i)
	}

	var received int
	var result int

	for digits := range ch {
		received++

		for _, d := range digits {
			if d == 1 || d == 4 || d == 7 || d == 8 {
				result++
			}
		}

		if received == len(inputs) {
			close(ch)
			break
		}
	}

	return strconv.Itoa(result)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	inputs := parseInput(r)
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

		if received == len(inputs) {
			close(ch)
			break
		}
	}

	return strconv.Itoa(sum)
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
