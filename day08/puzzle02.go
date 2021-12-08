package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type input struct {
	signalValues []signalValueInput
	outputValues []string
}

type signalValueInput struct {
	unknownValues []string
	knownValues   []string
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

	}

	return result
}

func parseSegments(input string) []string {
	input = strings.Trim(input, " \t\r\n")
	return strings.Split(input, " ")
}

func interpretSignalValues(inputs []string) []int {

	solvedNumbers := make(map[string]int)

	segmentMappings := make(map[rune]rune)

	unsolvedInputs := inputs

	tickIndex := 0

	for len(unsolvedInputs) > 0 {
		tickIndex++

		if tickIndex > 3 {
			break
		}

		unsolvedInputs = tick(unsolvedInputs, &solvedNumbers, &segmentMappings)

		fmt.Printf("*** TICK %d SOLVED: %v ***\n", tickIndex, solvedNumbers)
		fmt.Printf("*** TICK %d UNSOLVED: %v ***\n", tickIndex, unsolvedInputs)

	}

	fmt.Println(solvedNumbers)

	if len(solvedNumbers) != len(inputs) {
		panic(fmt.Sprintf("Expected %d solved inputs, but got %d", len(inputs), len(solvedNumbers)))
	}

	result := make([]int, 0, len(inputs))

	for i := 0; i < len(inputs); i++ {
		result = append(result, solvedNumbers[inputs[i]])
	}

	return result
}

func tick(inputSignalValues []string, solvedNumbers *map[string]int, segmentMappings *map[rune]rune) []string {
	segmentsForNumbers := []string{
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

	numbersForSegments := make(map[string]int)
	for number, segments := range segmentsForNumbers {
		numbersForSegments[segments] = number
	}

	// step 1: try and resolve any segments we can

	for _, inputSignalValue := range inputSignalValues {
		// inputSignalValue is e.g. "adf"

		// build a list of candidates that are the same length
		candidates := make(map[int]bool)
		var candidateNumber int
		for num := 0; num <= 9; num++ {

			isSolved := false
			for _, n := range *solvedNumbers {
				if n == num {
					isSolved = true
					break
				}
			}

			if isSolved {
				continue
			}

			if len(segmentsForNumbers[num]) == len(inputSignalValue) {
				candidates[num] = true
				candidateNumber = num
			}
		}

		if len(candidates) != 1 {
			continue
		}

		fmt.Printf("candidate for %s: %v\n", inputSignalValue, candidateNumber)

		candidateNumberSegments := segmentsForNumbers[candidateNumber]

		candidateMappings := make(map[rune][]rune)

		for _, inputSegment := range inputSignalValue {
			for _, knownSegment := range candidateNumberSegments {
				// if the # of input values w/ the input segment != # of known values w/ known segment, that's bad
				// also, if the individual input values w/ in the input segment don't have the same lengths as the known values, that's bad

				inputValuesWithSegment := findValuesWithSegment(inputSignalValues, inputSegment)
				sort.Slice(inputValuesWithSegment, func(i, j int) bool {
					return len(inputValuesWithSegment[i]) < len(inputValuesWithSegment[j])
				})

				knownValuesWithSegment := findValuesWithSegment(segmentsForNumbers, knownSegment)
				sort.Slice(knownValuesWithSegment, func(i, j int) bool {
					return len(knownValuesWithSegment[i]) < len(knownValuesWithSegment[j])
				})

				if len(inputValuesWithSegment) != len(knownValuesWithSegment) {
					continue
				}

				fmt.Printf("Input values with segment %s: %v\n", string(inputSegment), inputValuesWithSegment)
				fmt.Printf("Known values with segment %s: %v\n", string(knownSegment), knownValuesWithSegment)

				allLengthsMatch := true

				for i := range inputValuesWithSegment {
					allLengthsMatch = allLengthsMatch && len(inputValuesWithSegment[i]) == len(knownValuesWithSegment[i])
				}

				if !allLengthsMatch {
					continue
				}

				fmt.Printf("Candidate mapping: %s -> %s\n", string(inputSegment), string(knownSegment))
				candidateMappings[inputSegment] = append(candidateMappings[inputSegment], knownSegment)
			}
		}

		for inputSegment, candidateSegments := range candidateMappings {
			if len(candidateSegments) == 1 {
				fmt.Printf("segment %s == %s\n", string(inputSegment), string(candidateSegments[0]))
			}
			(*segmentMappings)[inputSegment] = candidateSegments[0]
		}
	}

	fmt.Printf("known segment mappings: %v\n", *segmentMappings)

	// step 2: given our known segment mappings, see if we can solve anything
	unsolvedInputs := make([]string, 0, len(inputSignalValues))
	potentialSolves := make(map[int][]string)

	for _, inputSignalValue := range inputSignalValues {

		// if every segment in this value can be mapped, do it
		ok := true
		segments := make([]rune, 0, len(inputSignalValue))

		for _, inputSegment := range inputSignalValue {
			knownSegment, found := (*segmentMappings)[inputSegment]
			fmt.Printf("%s = %s, %v\n", string(inputSegment), string(knownSegment), found)
			if !found {
				ok = false
				break
			}
			segments = append(segments, knownSegment)
		}

		if !ok {
			// try this input on the next tick
			fmt.Printf("Unsolvable input: %s\n", inputSignalValue)
			unsolvedInputs = append(unsolvedInputs, inputSignalValue)
			continue
		}

		knownValue := string(segments)
		number := numbersForSegments[knownValue]

		potentialSolves[number] = append(potentialSolves[number], inputSignalValue)

		fmt.Printf("____ %s maps to %s, which == %d ____\n", inputSignalValue, knownValue, number)
	}

	// step 3 for times where we were able to come up with 1 solution, keep it

	for number, inputSignalValues := range potentialSolves {

		fmt.Printf("%d solves %v\n", number, inputSignalValues)

		if len(inputSignalValues) > 1 {
			for _, i := range inputSignalValues {
				unsolvedInputs = append(unsolvedInputs, i)
			}
			continue
		}

		(*solvedNumbers)[inputSignalValues[0]] = number
	}

	return unsolvedInputs

}

func findValuesWithLength(values []string, length int) []string {
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

func removeSolvedInputs(input []signalValueInput) []signalValueInput {
	var result []signalValueInput
	for _, i := range input {
		if len(i.unknownValues) == 0 {
			result = append(result, i)
		}
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
