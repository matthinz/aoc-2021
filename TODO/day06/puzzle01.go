package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	days := 1

	if len(os.Args[1:]) > 0 {
		dayArg, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err != nil {
			panic(err)
		}
		days = int(dayArg)
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		numbers := parseNumbers(line)

		simulate(numbers, days)

	}

}

func simulate(numbers []int, days int) []int {

	const MaxTimer = 8
	timers := make([]int, MaxTimer+1)

	for _, timer := range numbers {
		timers[timer]++
	}

	for day := 1; day <= days; day++ {

		nextTick := make([]int, MaxTimer+1)

		for i := 0; i <= MaxTimer; i++ {
			count := timers[i]
			switch i {
			case 0:
				nextTick[6] += count
				nextTick[8] += count
				break
			default:
				nextTick[i-1] += count
				break
			}
		}

		timers = nextTick
		fmt.Printf("Day %d: %d lanternfish\n", day, sum(timers))

	}

	return numbers

}

func sum(values []int) int {
	var result int
	for _, value := range values {
		result += value
	}
	return result
}

func parseNumbers(input string) []int {
	var result []int
	for _, token := range strings.Split(input, ",") {
		value, err := strconv.ParseInt(token, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, int(value))
	}
	return result
}
