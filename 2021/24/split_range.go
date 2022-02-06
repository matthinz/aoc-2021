package d24

import "fmt"

type splitRelationship int

const (
	beforeSplit splitRelationship = iota
	atSplit
	afterSplit
)

// splitRange is a Range implementation that provides a subset of a parent
// range. Used to support .Split() implementations.
type splitRange struct {
	parent Range
	around int
	rel    splitRelationship
}

func newSplitRanges(parent Range, around int) (Range, Range, Range) {
	before := splitRange{parent, around, beforeSplit}
	at := splitRange{parent, around, atSplit}
	after := splitRange{parent, around, afterSplit}
	return &before, &at, &after
}

func (r *splitRange) Includes(value int) bool {
	for v := range r.Values() {
		if v == value {
			return true
		}
	}
	return false
}

func (r *splitRange) Split(around int) (Range, Range, Range) {
	return newSplitRanges(r, around)
}

func (r *splitRange) String() string {
	switch r.rel {
	case beforeSplit:
		return fmt.Sprintf("%s less than %d", r.parent.String(), r.around)
	case atSplit:
		return fmt.Sprintf("%s at %d", r.parent.String(), r.around)
	case afterSplit:
		return fmt.Sprintf("%s greater than %d", r.parent.String(), r.around)
	default:
		panic("Invalid relationship")
	}
}

func (r *splitRange) Values() chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for value := range r.parent.Values() {
			if r.rel == beforeSplit && value < r.around {
				ch <- value
			} else if r.rel == atSplit && value == r.around {
				ch <- value
			} else if r.rel == afterSplit && value > r.around {
				ch <- value
			}
		}
	}()
	return ch
}
