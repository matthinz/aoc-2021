package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {

	numbers := readBinaryNumberList(os.Stdin)

	gamma, epsilon := CalculateGammaAndEpsilon(numbers)

	fmt.Printf("gamma: %s (%d)\n", strconv.FormatInt(int64(gamma), 2), gamma)
	fmt.Printf("epsilon: %s (%d)\n", strconv.FormatInt(int64(epsilon), 2), epsilon)
	fmt.Printf("product: %d\n", gamma*epsilon)
}

func CalculateGammaAndEpsilon(numbers []int) (int, int) {
	bitLength := 16
	setBits := make([]int, bitLength)

	// Create a map of bit index to number of elements that have the corresponding bit set
	for _, value := range numbers {
		for i := 0; i < bitLength; i++ {
			mask := 1 << i
			bitIsSet := value&mask != 0
			if bitIsSet {
				setBits[i]++
			}
		}
	}

	fmt.Println(setBits)

	var gamma, epsilon int

	for i := 0; i < bitLength; i++ {
		countWithBitSet := setBits[i]

		if countWithBitSet == 0 {
			continue
		}

		if countWithBitSet >= len(numbers)/2 {
			// when >= 50% have bit set, set corresponding bit in gamma
			gamma = gamma | 1<<i
		} else {
			// when < 50% have bit set, set corresponding bit in epsilon
			epsilon = epsilon | 1<<i
		}
	}

	return gamma, epsilon
}

func readBinaryNumberList(r io.Reader) []int {
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
