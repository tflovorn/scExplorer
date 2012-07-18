// Functions related to traversal of the first Brillouin zone.
package bzone

import vec "../vector"

import (
	_ "fmt"
	"math"
)

type BzFunc func(k vec.Vector) float64
type bzConsumer func(next, total float64, other *float64) float64

// Sum values of fn over all Brillouin zone points. It must be safe to make
// multiple concurrent calls to fn.
// Uses Kahan summation algorithm for increased accuracy.
func Sum(pointsPerSide int, dimension int, fn BzFunc) float64 {
	add := func(next, total float64, c *float64) float64 {
		// add next to total; c holds error compensation information
		y := next - *c
		t := total + y
		*c = (t - total) - y
		return t
	}
	return bzReduce(add, 0.0, pointsPerSide, dimension, fn)
}

// Find the minimum of fn over all Brillouin zone points. It must be safe to
// make multiple concurrent calls to fn.
func Minimum(pointsPerSide int, dimension int, fn BzFunc) float64 {
	minimum := func(next, min float64, foo *float64) float64 {
		if next < min {
			return next
		}
		return min
	}
	return bzReduce(minimum, math.MaxFloat64, pointsPerSide, dimension, fn)
}

// Iterate over the Brillouin zone, accumulating the values of fn with combine.
// Dev note (TODO?): An interface for point generators could be created to
// allow for dependency injection here to replace bzPoints with an arbitrary
// point generator. By doing this we could treat symmetric functions
// differently (more efficiently). This would also allow testing bzReduce
// independently of bzPoints.
func bzReduce(combine bzConsumer, start float64, L, d int, fn BzFunc) float64 {
	maxWorkers := 2 // TODO: set this via config file
	points := bzPoints(L, d)
	work := func(result chan float64) {
		total := start
		other := 0.0
		for {
			k := <-points
			if k != nil {
				total = combine(fn(k), total, &other)
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
	other := 0.0
	for i := 0; i < maxWorkers; i++ {
		fullTotal = combine(<-results[i], fullTotal, &other)
	}
	return fullTotal
}

// Produce a channel whose values cover each Brillouin zone point once. 
// After all points have been traversed, the channel's values are nil.
// TODO: it would be nice to cache the result of this to avoid many
// re-generations for the same arguments.
func bzPoints(L, d int) <-chan vec.Vector {
	points := make(chan vec.Vector)
	// start is the minumum value of any component of a point
	start := -math.Pi
	// (finish - step) is the maximum value of any component of a point
	finish := -start
	// step is the separation between point components
	step := (finish - start) / float64(L)

	go func() {
		k := vec.ZeroVector(d)
		kIndex := make([]int, d)
		// set initial value for k
		for i := 0; i < d; i++ {
			k[i] = start
			kIndex[i] = 0
		}
		// iterate over Brillouin zone
		done := false
		for {
			if !done {
				//fmt.Println(k)
				points <- k
			} else {
				break
			}
			done = bzAdvance(k, kIndex, start, step, L, d)
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
func bzAdvance(k vec.Vector, kIndex []int, start, step float64, L, d int) bool {
	for i := 0; i < d; i++ {
		// Check if we need to carry.
		// To avoid any possible (maybe imagined) risk of
		// rounding error breaking results, use kIndex instead
		// of making this comparison as k[i] == finish.
		carry := kIndex[i] == L
		if !carry {
			// Advance k.
			k[i] += step
			kIndex[i] += 1
			// Break unless advancing k caused us to need
			// to carry.
			if kIndex[i] != L {
				break
			}
		}
		// If we get to here, a carry is required.
		if i == d-1 {
			// finished iteration
			return true
		}
		k[i] = start
		kIndex[i] = 0
	}
	return false
}
