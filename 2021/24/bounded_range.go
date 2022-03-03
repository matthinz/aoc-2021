package d24

// BoundedRange defines a range with known (inclusive) minimum and maximum values.
type BoundedRange interface {
	Range
	Min() int
	Max() int
}
