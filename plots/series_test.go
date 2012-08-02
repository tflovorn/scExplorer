package plots

import (
	"sort"
	"testing"
)

type seriesTestData struct {
	X, Y, Z float64
}

func seriesTestDefaultData(n int) []interface{} {
	data := make([]seriesTestData, n)
	ints := make([]interface{}, n)
	for i := 0; i < n; i++ {
		ifl, nfl := float64(i), float64(n)
		data[i] = seriesTestData{nfl - ifl, ifl, float64(i % 2)}
		ints[i] = data[i]
	}
	return ints
}

// extractXY should return one series with its x values sorted.
func TestExtractXY(t *testing.T) {
	vals := seriesTestDefaultData(5)
	series := extractXY(vals, "X", "Y")
	if len(series) != 1 {
		t.Fatalf("extractXY returned > 1 series")
	}
	s := series[0]
	if !sort.Float64sAreSorted(s.xs) {
		t.Fatalf("extractXY returned incorrectly sorted series")
	}
}

// ExtractSeries should return a []Series with length equal to the number of
// values of z passed in to it. Individual series should have their x values
// sorted.
func TestExtractSeries(t *testing.T) {
	vals := seriesTestDefaultData(6)
	series := ExtractSeries(vals, []string{"X", "Y", "Z"})
	expectedLength := 2
	if len(series) != expectedLength {
		t.Fatalf("ExtractSeries returned incorrect number of series")
	}
	for i := 0; i < expectedLength; i++ {
		s := series[i]
		if !sort.Float64sAreSorted(s.xs) {
			t.Fatalf("ExtractSeries returned incorrectly sorted series")
		}
	}
}
