package d24

type emptyRange struct {
}

var EmptyRange = &emptyRange{}

func (r *emptyRange) Includes(value int) bool {
	return false
}

func (r *emptyRange) Split() []ContinuousRange {
	return []ContinuousRange{}
}

func (r *emptyRange) String() string {
	return "<>"
}

// Returns a function that, when executed, returns the next value in the
// range along with a boolean indicating whether the call succeeded.
func (r *emptyRange) Values(context string) func() (int, bool) {
	return func() (int, bool) {
		return 0, false
	}
}
