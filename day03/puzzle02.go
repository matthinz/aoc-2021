package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {

	numbers := readNumbers(os.Stdin)

	o2Rating := FindO2GeneratorRating(numbers)
	co2ScrubberRating := FindCo2ScrubberRating(numbers)

	fmt.Printf("Total numbers: %d\n", len(numbers))
	fmt.Printf("O2 generator rating: %d\n", o2Rating)
	fmt.Printf("CO2 scrubber rating: %d\n", co2ScrubberRating)
	fmt.Printf("Life support rating: %d\n", o2Rating*co2ScrubberRating)
}

func FindCo2ScrubberRating(numbers []string) int64 {
	candidates := numbers
	pos := 0

	for {
		_, leastCommonBit := FindMostAndLeastCommonBits(candidates, pos)
		candidates = FilterNumbersByBitAtPosition(candidates, leastCommonBit, pos)

		switch len(candidates) {
		case 0:
			panic("Ran out of candidates!")
		case 1:
			result, err := strconv.ParseInt(candidates[0], 2, 64)
			if err != nil {
				panic(err)
			}
			return result
		}

		pos++
	}
}

func FindO2GeneratorRating(numbers []string) int64 {
	candidates := numbers
	pos := 0

	for {

		mostCommonBit, _ := FindMostAndLeastCommonBits(candidates, pos)
		candidates = FilterNumbersByBitAtPosition(candidates, mostCommonBit, pos)

		switch len(candidates) {
		case 0:
			panic("Ran out of candidates!")
		case 1:
			result, err := strconv.ParseInt(candidates[0], 2, 64)
			if err != nil {
				panic(err)
			}
			return result
		}

		pos++
	}

}

func FilterNumbersByBitAtPosition(numbers []string, bit int, pos int) []string {
	result := make([]string, 0, len(numbers))

	for _, value := range numbers {

		isMatch := bit == 1 && value[pos] == '1' || bit == 0 && value[pos] == '0'
		if isMatch {
			result = append(result, value)
		}
	}

	return result
}

func FindMostAndLeastCommonBits(numbers []string, position int) (int, int) {
	var oneCount int

	for _, value := range numbers {
		if value[position] == '1' {
			oneCount++
		}
	}

	half := float64(len(numbers)) / 2.0

	if float64(oneCount) == half {
		// when equally common, we say "1" for most common and "0" for least common
		return 1, 0
	}

	var mostCommon int
	var leastCommon int

	if float64(oneCount) >= half {
		// when >= 50% are 1, we call 1 the most common
		mostCommon = 1
	}

	if float64(oneCount) < half {
		// when < 50% are 1, we call 1 the least common
		leastCommon = 1
	}

	return mostCommon, leastCommon

}

func readNumbers(r io.Reader) []string {
	scanner := bufio.NewScanner(os.Stdin)
	var bitLength int
	var numbers []string

	for scanner.Scan() {
		value := scanner.Text()
		if len(value) == 0 {
			continue
		}

		if bitLength == 0 {
			bitLength = len(value)
		} else if len(value) != bitLength {
			panic("bad bitLength")
		}
		numbers = append(numbers, value)
	}

	return numbers
}
