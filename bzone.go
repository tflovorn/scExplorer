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
		// add next to total; c holds error compensation information
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

// Iterate over the Brillouin zone, accumulating the values of fn with combine.
func bzReduce(combine bzConsumer, start float64, pointsPerSide int64, dimension int64, fn BzFunc) float64 {
	maxWorkers := 2 // TODO: set this via config file
	points, done := bzPoints(pointsPerSide, dimension)
	work := func(result chan float64) {
		var k BzPoint
		total := start
		for {
			select {
			case k = <-points:
				combine(fn(k), &total)
			case <-done:
				done <- true
				break
			}
		}
		result <- total
	}
	// set up workers
	results := make([]chan float64, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		results[i] = make(chan float64)
		go work(results[i])
	}
	// collect the results
	fullTotal := start
	for i := 0; i < maxWorkers; i++ {
		combine(<-results[i], &fullTotal)
	}
	return fullTotal
}

// Produce (points, done) where points is a channel whose values cover each
// Brillouin zone point once, and done is is a channel which contains true
// after all points have been traversed.
func bzPoints(pointsPerSide int64, dimension int64) (chan BzPoint, chan bool) {
	points := make(chan BzPoint)
	done := make(chan bool)
	// start is the minumum value of any component of a point
	start := -math.Pi
	// (finish - step) is the maximum value of any component of a point
	finish := -start
	// separation between point components
	step := (finish - start) / float64(pointsPerSide)

	generatePoints := func() {
		k := make([]float64, dimension)
		for i := int64(0); i < dimension; i++ {
			k[i] = start
		}
		for {
			// TODO: iterate over points
			k[0] += step
			break
		}
		done <- true
	}
	go generatePoints()
	return points, done
}
