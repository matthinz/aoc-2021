package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var bitLength int
	var numberCount int
	// tracks the number of '1' bits at each corresponding position
	var oneCounts []int

	for scanner.Scan() {
		value := scanner.Text()
		if len(value) == 0 {
			continue
		}

		if bitLength == 0 {
			bitLength = len(value)
			oneCounts = make([]int, bitLength)
		} else if len(value) != bitLength {
			panic("invalid bit length")
		}

		numberCount++

		for index, c := range value {
			if c == '1' {
				oneCounts[index]++
			}
		}
	}

	// gamma is a `bitLength`-length binary number whe
	var gamma uint
	var epsilon uint

	for index, count := range oneCounts {
		if count >= numberCount/2 {
			// this bit should be 1
			gamma = gamma | 1<<(bitLength-(index)-1)
		} else {
			epsilon = epsilon | 1<<(bitLength-(index)-1)
		}
	}

	fmt.Printf("gamma: %s (%d)\n", strconv.FormatInt(int64(gamma), 2), gamma)
	fmt.Printf("epsilon: %s (%d)\n", strconv.FormatInt(int64(epsilon), 2), epsilon)
	fmt.Printf("product: %d\n", gamma*epsilon)

}
