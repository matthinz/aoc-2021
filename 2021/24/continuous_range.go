package d24

import "fmt"

type continuousRange struct {
	min, max int
	step     int
}

func newContinuousRange(min, max, step int) *continuousRange {
	if max < min {
		temp := max
		max = min
		min = temp
	}
	return &continuousRange{min, max, step}
}

func (r *continuousRange) Includes(value int) bool {
	return value >= r.min && value <= r.max
}

func (r *continuousRange) Split(around int) (Range, Range, Range) {

	if around < r.min {
		return nil, nil, r
	}

	if around > r.max {
		return r, nil, nil
	}

	before := splitRange{r, around, beforeSplit}
	after := splitRange{r, around, afterSplit}

	var at Range

	if around-r.min%r.step == 0 {
		at = &splitRange{r, around, atSplit}
	} else {
		at = nil
	}

	return &before, at, &after
}

func (r *continuousRange) String() string {
	if r.min == r.max {
		return fmt.Sprintf("%d", r.min)
	} else if r.step != 1 {
		return fmt.Sprintf("%d..%d step %d", r.min, r.max, r.step)
	} else {
		return fmt.Sprintf("%d..%d", r.min, r.max)
	}
}

func (r *continuousRange) Values() chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for value := r.min; value <= r.max; value += r.step {
			ch <- value
		}
	}()

	return ch
}
