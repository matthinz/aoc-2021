package d02

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

type policy struct {
	char rune
	a    int
	b    int
}

type input struct {
	policy   policy
	password string
}

func New() aoc.Day {
	return aoc.NewDay(2, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	inputs := parseInput(r)
	valid := 0
	for _, i := range inputs {
		if i.isValidForSledRentalPlace() {
			valid++
		}
	}
	return strconv.Itoa(valid)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	inputs := parseInput(r)
	valid := 0
	for _, i := range inputs {
		if i.isValidForOTCA() {
			valid++
		}
	}

	return strconv.Itoa(valid)
}

func (i *input) isValidForOTCA() bool {

	aChar := rune(i.password[i.policy.a-1])
	bChar := rune(i.password[i.policy.b-1])

	// Exactly one of these must match the policy char
	return ((aChar == i.policy.char) && (bChar != i.policy.char)) || ((bChar == i.policy.char) && aChar != i.policy.char)
}

func (i *input) isValidForSledRentalPlace() bool {
	count := 0
	for _, r := range i.password {
		if r == i.policy.char {
			count++
		}
	}

	min := i.policy.a
	max := i.policy.b

	return count >= min && count <= max
}

func parseInput(r io.Reader) []input {
	result := make([]input, 0)
	rx := regexp.MustCompile("^(\\d+)-(\\d+) (\\w): (.+)$")

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		m := rx.FindStringSubmatch(line)
		if m == nil {
			panic(fmt.Sprintf("Not matched: %s", line))
		}

		min, _ := strconv.ParseInt(m[1], 10, 32)
		max, _ := strconv.ParseInt(m[2], 10, 32)

		result = append(result, input{
			password: m[4],
			policy: policy{
				char: rune(m[3][0]),
				a:    int(min),
				b:    int(max),
			},
		})
	}

	return result
}
