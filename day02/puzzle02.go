package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var x int64
	var y int64
	var aim int64

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			continue
		}
		command := tokens[0]
		value, err := strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			panic(err)
		}

		switch command {
		case "forward":
			x += value
			y += aim * value
			break
		case "down":
			aim += value
			break
		case "up":
			aim -= value
			break
		}

		fmt.Printf("%s %d: %d x %d = %d\n", command, value, x, y, x*y)
	}
}
