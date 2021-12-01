package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const WindowSize = 3

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var window []int64
	var prevSum *int64
	var increases int64

	for scanner.Scan() {
		value, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			continue
		}
		window = append(window, value)

		if len(window) == WindowSize {

			sum := sumValues(window)

			if prevSum != nil {
				increased := sum > *prevSum
				decreased := sum < *prevSum

				if increased {
					fmt.Printf("%s = %d (increased)\n", formatEquation(window), sum)
					increases++
				} else if decreased {
					fmt.Printf("%s = %d (decreased)\n", formatEquation(window), sum)
				}
			} else {
				fmt.Printf("%s = %d (n/a)\n", formatEquation(window), sum)
			}

			prevSum = &sum

			window = window[1:]
		}
	}

	fmt.Printf("\n\n\n%d increases\n", increases)
}

func sumValues(values []int64) int64 {
	var result int64
	for _, v := range values {
		result += v
	}
	return result
}

func formatEquation(values []int64) string {
	builder := strings.Builder{}
	for index, v := range values {
		if index > 0 {
			builder.WriteString("+")
		}
		builder.WriteString(strconv.FormatInt(v, 10))
	}
	return builder.String()
}
