package d24

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed input
var realInput string

func TestParseRealInput(t *testing.T) {
	parseInput(strings.NewReader(realInput))
}
