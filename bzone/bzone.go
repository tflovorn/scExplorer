// Functions related to traversal of the first Brillouin zone.
package bzone

import "math"

type BzPoint []float64
type BzFunc func(k BzPoint) float64
type bzConsumer func(next float64, total *float64)

// Sum values of fn over all Brillouin zone points.
// Uses Kahan summation algorithm for increased accuracy.
func Sum(pointsPerSide int, dimension int, fn BzFunc) float64 {
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
func Minimum(pointsPerSide int, dimension int, fn BzFunc) float64 {
	minimum := func(next float64, min *float64) {
		if next < *min {
			*min = next
		}
	}
	return bzReduce(minimum, math.MaxFloat64, pointsPerSide, dimension, fn)
}

// Iterate over the Brillouin zone, accumulating the values of fn with combine.
// Dev note (TODO?): An interface for point generators could be created to
// allow for dependency injection here to replace bzPoints with an arbitrary
// point generator. By doing this we could treat symmetric functions
// differently (more efficiently). This would also allow testing bzReduce
// independently of bzPoints.
func bzReduce(combine bzConsumer, start float64, pointsPerSide int, dimension int, fn BzFunc) float64 {
	maxWorkers := 2 // TODO: set this via config file
	points := bzPoints(pointsPerSide, dimension)
	work := func(result chan float64) {
		var k BzPoint
		total := start
		for {
			k = <-points
			if k != nil {
				combine(fn(k), &total)
			} else {
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

// Produce a channel whose values cover each Brillouin zone point once. 
// After all points have been traversed, the channel's values are nil.
// TODO: it would be nice to cache the result of this to avoid many
// re-generations for the same arguments.
func bzPoints(pointsPerSide int, dimension int) <-chan BzPoint {
	points := make(chan BzPoint)
	// start is the minumum value of any component of a point
	start := -math.Pi
	// (finish - step) is the maximum value of any component of a point
	finish := -start
	// step is the separation between point components
	step := (finish - start) / float64(pointsPerSide)

	go func() {
		k := make([]float64, dimension)
		kIndex := make([]int, dimension)
		// set initial value for k
		for i := 0; i < dimension; i++ {
			k[i] = start
			kIndex[i] = 0
		}
		// iterate over Brillouin zone
		done := false
		for {
			if !done {
				points <- k
			} else {
				break
			}
			done = bzAdvance(k, kIndex, start, step, pointsPerSide, dimension)
		}
		// we're done; signal that
		for {
			points <- nil
		}
	}()
	return points
}

// Advances k and kIndex to the next value. Returns true on overflow back to
// the initial k value.
func bzAdvance(k []float64, kIndex []int, start, step float64, pointsPerSide, dimension int) bool {
	for i := 0; i < dimension; i++ {
		// Check if we need to carry.
		// To avoid any possible (maybe imagined) risk of
		// rounding error breaking results, use kIndex instead
		// of making this comparison as k[i] == finish.
		carry := kIndex[i] == pointsPerSide
		if !carry {
			// Advance k.
			k[i] += step
			kIndex[i] += 1
			// Break unless advancing k caused us to need
			// to carry.
			if kIndex[i] != pointsPerSide {
				break
			}
		}
		// If we get to here, a carry is required.
		if i == dimension-1 {
			// finished iteration
			return true
		}
		k[i] = start
		kIndex[i] = 0
	}
	return false
}
