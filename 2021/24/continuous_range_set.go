package d24

import (
	"fmt"
	"sort"
	"strings"
)

type continuousRangeSet struct {
	members []ContinuousRange
}

func NewContinuousRangeSet(ranges []ContinuousRange) Range {
	members := make([]ContinuousRange, len(ranges))
	copy(members, ranges)

	sort.Slice(members, func(i, j int) bool {
		if members[i].Min() < members[j].Min() {
			return true
		} else if members[i].Min() == members[j].Min() && members[i].Max() < members[j].Max() {
			return true
		} else {
			return false
		}
	})

	return &continuousRangeSet{members}
}

func (c *continuousRangeSet) Includes(value int) bool {
	for _, m := range c.members {
		if m.Includes(value) {
			return true
		}
	}
	return false
}

func (c *continuousRangeSet) Split() []ContinuousRange {
	return c.members
}

func (c *continuousRangeSet) Values(context string) func() (int, bool) {
	var next func() (int, bool)
	index := 0

	return func() (int, bool) {
		for {
			if next == nil {
				if index >= len(c.members) {
					return 0, false
				}
				next = c.members[index].Values(context)
				index++
			}

			value, ok := next()

			if ok {
				return value, ok
			} else {
				next = nil
			}
		}
	}
}

func (c *continuousRangeSet) String() string {
	if len(c.members) == 1 {
		return c.members[0].String()
	}

	items := make([]string, len(c.members))
	for i, m := range c.members {
		items[i] = m.String()
	}
	return fmt.Sprintf("<%s>", strings.Join(items, ","))
}
