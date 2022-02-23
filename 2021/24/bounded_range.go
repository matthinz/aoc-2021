package d24

// BoundedRange defines a range of values with a known minimum and maximum.
type BoundedRange interface {
	Min() int
	Max() int
}

type boundedRange struct {
	min, max int
}

func (r *boundedRange) Max() int {
	return r.max
}

func (r *boundedRange) Min() int {
	return r.min
}

// Attempts to read the bounds of the given range.
func getBounds(r Range) BoundedRange {

	if b, isBounded := r.(BoundedRange); isBounded {
		return b
	}

	if compound, isCompound := r.(*compoundRange); isCompound {

		allHaveBounds := true
		var min, max *int
		for _, m := range compound.members {
			b := getBounds(m)
			if b == nil {
				allHaveBounds = false
				break
			} else {
				bMin := b.Min()
				bMax := b.Max()
				if min == nil || bMin < *min {
					min = &bMin
				}
				if max == nil || bMax > *max {
					max = &bMax
				}
			}
		}

		if !allHaveBounds || min == nil || max == nil {
			return nil
		}

		return &boundedRange{*min, *max}
	}

	return nil
}
