package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var prevValue *int64
	var increases int64

	for scanner.Scan() {
		value, err := strconv.ParseInt(scanner.Text(), 10, 32)

		if err != nil {
			continue
		}

		if prevValue != nil {

			increased := value > *prevValue
			decreased := value < *prevValue

			if increased {
				increases++
				fmt.Printf("%d: increased\n", value)
			} else if decreased {
				fmt.Printf("%d: decreased\n", value)
			}

		}

		prevValue = &value
	}

	fmt.Printf("\n\n\n%d increases\n", increases)
}
