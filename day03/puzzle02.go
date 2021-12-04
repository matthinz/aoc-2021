package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const MaxBitLength = 16

func main() {

	numbers := readNumbers(os.Stdin)

	o2Rating := FindO2GeneratorRating(numbers)
	co2ScrubberRating := FindCo2ScrubberRating(numbers)

	fmt.Printf("Total numbers: %d\n", len(numbers))
	fmt.Printf("O2 generator rating: %d\n", o2Rating)
	fmt.Printf("CO2 scrubber rating: %d\n", co2ScrubberRating)
	fmt.Printf("Life support rating: %d\n", o2Rating*co2ScrubberRating)
}

func FindCo2ScrubberRating(numbers []int) int {
	candidates := numbers
	pos := MaxBitLength

	for {
		_, leastCommonBit := FindMostAndLeastCommonBits(candidates, pos)
		candidates = FilterNumbersByBitAtPosition(candidates, leastCommonBit, pos)

		switch len(candidates) {
		case 0:
			panic("Ran out of candidates!")
		case 1:
			return candidates[0]
		}

		pos--
	}
}

func FindO2GeneratorRating(numbers []int) int {
	candidates := numbers
	pos := MaxBitLength

	for {

		mostCommonBit, _ := FindMostAndLeastCommonBits(candidates, pos)
		candidates = FilterNumbersByBitAtPosition(candidates, mostCommonBit, pos)

		switch len(candidates) {
		case 0:
			panic("Ran out of candidates!")
		case 1:
			return candidates[0]
		}

		pos--
	}

}

func FilterNumbersByBitAtPosition(numbers []int, bit int, pos int) []int {
	result := make([]int, 0, len(numbers))
	mask := 1 << pos

	for _, value := range numbers {

		var isMatch bool

		if bit == 0 {
			// check the bit is _not set_ at the position
			isMatch = value&mask == 0
		} else {
			// check the bit is set at the p
			isMatch = value&mask != 0
		}

		if isMatch {
			result = append(result, value)
		}
	}

	return result
}

func FindMostAndLeastCommonBits(numbers []int, position int) (int, int) {

	mask := 1 << position
	setCount := 0

	for _, value := range numbers {
		bitIsSet := value&mask != 0
		if bitIsSet {
			setCount++
		}
	}

	if setCount == 0 {
		// no numbers have this bit set, so just ignore
		return 0, 0
	}

	half := float64(len(numbers)) / 2.0

	if float64(setCount) == half {
		// when equally common, we say "1" for most common and "0" for least common
		return 1, 0
	}

	var mostCommon int
	var leastCommon int

	if float64(setCount) >= half {
		// when >= 50% are 1, we call 1 the most common
		mostCommon = 1
	}

	if float64(setCount) < half {
		// when < 50% are 1, we call 1 the least common
		leastCommon = 1
	}

	return mostCommon, leastCommon
}

func readNumbers(r io.Reader) []int {
	scanner := bufio.NewScanner(os.Stdin)
	var result []int

	for scanner.Scan() {
		token := scanner.Text()
		if len(token) == 0 {
			continue
		}

		value, err := strconv.ParseInt(token, 2, 64)

		if err != nil {
			panic(err)
		}

		result = append(result, int(value))
	}

	return result
}
