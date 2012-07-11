// Functions related to traversal of the first Brillouin zone.
package scExplorer

import "math"

type BzFunc func(k []float64) float64
type consumer func(next float64, total *float64)

// Sum values of fn over all Brillouin zone points.
// Uses Kahan summation algorithm for increased accuracy.
func BzSum(pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	c := 0.0
	var y, t float64
	add := func(next float64, total *float64) {
		y = next - c
		t = *total + y
		c = (t - *total) - y
		*total = t
	}
	return reduce(add, 0.0, pointsPerSide, dimension, fn)
}

// Find the minimum of fn over all Brillouin zone points.
func BzMinimum(pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	minimum := func(next float64, min *float64) {
		if next < *min {
			*min = next
		}
	}
	return reduce(minimum, math.MaxFloat64, pointsPerSide, dimension, fn)
}

func reduce(cs consumer, start float64, pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	return 0.0
}

func bzPoints(pointsPerSide int64, dimension int64) (chan []float64, chan bool) {
	return nil, nil
}
