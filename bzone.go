// Functions related to traversal of the first Brillouin zone.
package scExplorer

import "math"

type BzPoint []float64
type BzFunc func(k BzPoint) float64
type bzConsumer func(next float64, total *float64)

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
	return bzReduce(add, 0.0, pointsPerSide, dimension, fn)
}

// Find the minimum of fn over all Brillouin zone points.
func BzMinimum(pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	minimum := func(next float64, min *float64) {
		if next < *min {
			*min = next
		}
	}
	return bzReduce(minimum, math.MaxFloat64, pointsPerSide, dimension, fn)
}

// Iterate over the Brillouin zone, accumulating the values of fn with cs.
func bzReduce(cs bzConsumer, start float64, pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	maxWorkers := 2 // TODO: set this via config file
	points, done := bzPoints(pointsPerSide, dimension)
	results := make([]chan float64, maxWorkers)
	work := func(result chan float64) {
		var k BzPoint
		total := start
		for {
			select {
			case k = <-points:
				cs(fn(k), &total)
			case <-done:
				done <- true
				break
			}
		}
		result <- total
	}

	for i := 0; i < maxWorkers; i++ {
		results[i] = make(chan float64)
		go work(results[i])
	}
	fullTotal := start
	for i := 0; i < maxWorkers; i++ {
		cs(<-results[i], &fullTotal)
	}
	return fullTotal
}

// Produce (points, done) where points is a channel whose values cover each
// Brillouin zone point once, and done is is a channel which contains true
// after all points have been traversed.
func bzPoints(pointsPerSide int64, dimension int64) (chan BzPoint, chan bool) {
	return nil, nil
}
