package d24

type emptyRange struct {
}

var EmptyRange = &emptyRange{}

func (r *emptyRange) Includes(value int) bool {
	return false
}

// Splits this range around <around>. Returns 3 new Ranges representing
// the portion of the range less than <around>, <around> itself, and the
// portion of the range greater than <around>.
func (r *emptyRange) Split(around Range) (Range, Range, Range) {
	return nil, nil, nil
}

func (r *emptyRange) String() string {
	return "<>"
}

// Returns a function that, when executed, returns the next value in the
// range along with a boolean indicating whether the call succeeded.
func (r *emptyRange) Values() func() (int, bool) {
	return func() (int, bool) {
		return 0, false
	}
}
