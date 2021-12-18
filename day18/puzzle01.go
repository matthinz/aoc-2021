package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type snailfishNumberKind int

const (
	regularNumberKind = snailfishNumberKind(0)
	pairNumberKind    = iota
)

type snailfishNumber struct {
	kind   snailfishNumberKind
	value  int
	left   *snailfishNumber
	right  *snailfishNumber
	parent *snailfishNumber
}

func main() {

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		num, err := parseSnailfishNumber(line)
		if err != nil {
			panic(err)
		}

		fmt.Println(num.String())

	}

}

func parseSnailfishNumber(input string) (*snailfishNumber, error) {

	var n *snailfishNumber

	for pos, r := range input {

		switch {

		case r == '[':
			// open a new pair
			n = &snailfishNumber{
				kind:   pairNumberKind,
				parent: n,
			}

			if n.parent != nil {
				if n.parent.left == nil {
					n.parent.left = n
				} else if n.parent.right == nil {
					n.parent.right = n
				} else {
					return nil, fmt.Errorf("invalid")
				}
			}

		case r == ',':

		case r == ']':
			// close a pair

			if n.parent == nil {
				return n, nil
			} else {
				n = n.parent
			}

		// NOTE: Assuming number will always be < 10
		case r >= '0' && r <= '9':
			value, _ := strconv.Atoi(string(r))

			reg := snailfishNumber{
				kind:   regularNumberKind,
				value:  value,
				parent: n,
			}

			if n.left == nil {
				n.left = &reg
			} else if n.right == nil {
				n.right = &reg
			} else {
				panic("neither left nor right is nil")
			}

		default:
			return nil, fmt.Errorf("Invalid character found at %d: '%s'", pos, string(r))

		}

	}

	return nil, fmt.Errorf("Parsing failed")
}

////////////////////////////////////////////////////////////////////////////////
// snailfishNumber impl

func (s *snailfishNumber) String() string {
	switch s.kind {
	case regularNumberKind:
		return strconv.Itoa(s.value)
	case pairNumberKind:
		left := "nil"
		right := "nil"
		if s.left != nil {
			left = s.left.String()
		}
		if s.right != nil {
			right = s.right.String()
		}
		return fmt.Sprintf("[%s,%s]", left, right)
	default:
		panic("Invalid kind")
	}
}
