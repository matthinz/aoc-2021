package main

import (
	"bufio"
	"fmt"
	"math"
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

	var num *snailfishNumber

	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		parsed, err := parseSnailfishNumber(line)

		if err != nil {
			panic(err)
		}

		if num == nil {
			num = parsed
		} else {
			num = num.add(parsed)
		}
	}

	fmt.Println(num.magnitude())

}

// add combines s with other and reduces the result
// it returns a new snailfishNumber
func (s *snailfishNumber) add(other *snailfishNumber) *snailfishNumber {

	var left, right *snailfishNumber

	if s != nil {
		left = s.clone()
	}

	if other != nil {
		right = other.clone()
	}

	result := snailfishNumber{
		kind:  pairNumberKind,
		left:  left,
		right: right,
	}

	if result.left != nil {
		result.left.parent = &result
	}

	if result.right != nil {
		result.right.parent = &result
	}

	result.reduce()

	return &result
}

// clones this number. the clone will not have a parent set
func (s *snailfishNumber) clone() *snailfishNumber {

	if s == nil {
		return nil
	}

	c := snailfishNumber{
		kind:  s.kind,
		value: s.value,
		left:  s.left.clone(),
		right: s.right.clone(),
	}

	if c.left != nil {
		c.left.parent = &c
	}
	if c.right != nil {
		c.right.parent = &c
	}

	return &c

}

// depth returns the # of ancestors a number has
func (s *snailfishNumber) depth() int {
	if s.parent == nil {
		return 0
	}
	return s.parent.depth() + 1
}

func (s *snailfishNumber) explode() {

	// To explode a pair, the pair's left value is added to the first regular
	// number to the left of the exploding pair (if any), and the pair's right
	// value is added to the first regular number to the right of the exploding
	// pair (if any). Exploding pairs will always consist of two regular numbers.
	// Then, the entire exploding pair is replaced with the regular number 0.

	if s.kind != pairNumberKind {
		return
	}

	haveTwoRegularNumbers := s.left.kind == regularNumberKind && s.right.kind == regularNumberKind

	if !haveTwoRegularNumbers {
		panic("should have two regular numbers")
	}

	// Add the left value to the first number to the left

	var numberToTheLeft *snailfishNumber

	for n := s; n != nil; n = n.parent {
		isRight := n.parent != nil && n.parent.right == n
		if !isRight {
			continue
		}
		numberToTheLeft = n.parent.left
		break
	}

	if numberToTheLeft != nil {
		// find the rightmost literal inside numberToTheLeft
		for numberToTheLeft.kind != regularNumberKind {
			numberToTheLeft = numberToTheLeft.right
		}

		numberToTheLeft.value += s.left.value
	}

	// Add the right value to the first number to the right
	var numberToTheRight *snailfishNumber
	for n := s; n != nil; n = n.parent {
		isLeft := n.parent != nil && n.parent.left == n
		if !isLeft {
			continue
		}
		numberToTheRight = n.parent.right
		break
	}

	if numberToTheRight != nil {
		// find the leftmost literal inside numberToTheRight
		for numberToTheRight.kind != regularNumberKind {
			numberToTheRight = numberToTheRight.left
		}

		numberToTheRight.value += s.right.value
	}

	// Reset the pair to the regular number "0"

	s.kind = regularNumberKind
	s.value = 0
	s.left = nil
	s.right = nil

}

func (s *snailfishNumber) findLeftmost(f func(*snailfishNumber) bool) *snailfishNumber {
	if f(s) {
		return s
	}

	if s.kind != pairNumberKind {
		return nil
	}

	result := s.left.findLeftmost(f)
	if result != nil {
		return result
	}

	return s.right.findLeftmost(f)
}

func (s *snailfishNumber) findRightmost(f func(*snailfishNumber) bool) *snailfishNumber {
	if f(s) {
		return s
	}
	if s.kind != pairNumberKind {
		return nil
	}

	result := s.right.findRightmost(f)
	if result != nil {
		return result
	}

	return s.left.findRightmost(f)

}

func (s *snailfishNumber) magnitude() int {
	switch s.kind {
	case regularNumberKind:
		return s.value

	case pairNumberKind:
		return (3 * s.left.magnitude()) + (2 * s.right.magnitude())

	default:
		panic("invalid kind")
	}
}

func (s *snailfishNumber) checkDepths(expected int) {
	actual := s.depth()
	if actual != expected {
		panic(fmt.Sprintf("Depth for %s is wrong: expected %d, got %d", s.String(), expected, actual))
	}
	if s.kind == pairNumberKind {
		s.left.checkDepths(expected + 1)
		s.right.checkDepths(expected + 1)
	}
}

func (s *snailfishNumber) reduce() {
	s.reduceWithDebugging(false)
}

func (s *snailfishNumber) reduceWithDebugging(debug bool) {

	if debug {
		fmt.Printf("reducing: %s\n", s.String())
	}

	for {

		leftmostNested := s.findLeftmost(func(n *snailfishNumber) bool {

			// fmt.Printf("visit: %s (%d)\n", n.String(), n.depth())

			if n.kind != pairNumberKind {
				return false
			}

			// it has to be *explodable*
			if n.left.kind != regularNumberKind || n.right.kind != regularNumberKind {
				return false
			}

			return n.depth() >= 4
		})

		if leftmostNested != nil {

			if debug {
				value := leftmostNested.String()
				before := s.String()

				leftmostNested.explode()

				fmt.Printf("reduce(): explode %s %s -> %s\n", value, before, s.String())
			} else {
				leftmostNested.explode()
			}

			continue
		}

		leftmost10OrGreater := s.findLeftmost(func(n *snailfishNumber) bool {
			return n.kind == regularNumberKind && n.value >= 10
		})

		if leftmost10OrGreater != nil {
			if debug {
				value := leftmost10OrGreater.value
				before := s.String()
				leftmost10OrGreater.split()
				fmt.Printf("reduce(): split %d %s -> %s\n", value, before, s.String())
			} else {
				leftmost10OrGreater.split()
			}
			continue
		}

		if debug {
			fmt.Printf("reduce(): result %s\n", s.String())
		}

		return
	}
}

func (s *snailfishNumber) split() {
	if s.kind != regularNumberKind {
		return
	}

	left := snailfishNumber{
		kind:   regularNumberKind,
		value:  s.value / 2,
		parent: s,
	}

	right := snailfishNumber{
		kind:   regularNumberKind,
		value:  int(math.Ceil(float64(s.value) / 2.0)),
		parent: s,
	}

	s.kind = pairNumberKind
	s.value = 0
	s.left = &left
	s.right = &right
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
