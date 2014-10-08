// Functions related to traversal of the first Brillouin zone of a square
// lattice.
package bzone

import (
	"math"
)
import vec "github.com/tflovorn/scExplorer/vector"

type BzFunc func(k vec.Vector) float64
type bzConsumer func(next, total float64) float64

// Sum values of fn over all Brillouin zone points.
// Uses Kahan summation algorithm for increased accuracy.
func Sum(pointsPerSide int, dimension int, fn BzFunc) float64 {
	c := 0.0
	add := func(next, total float64) float64 {
		// add next to total; c holds error compensation information
		y := next - c
		t := total + y
		c = (t - total) - y
		return t
	}
	return bzReduce(add, 0.0, pointsPerSide, dimension, fn)
}

// Average = Sum / (total number of points)
func Avg(pointsPerSide int, dimension int, fn BzFunc) float64 {
	N := math.Pow(float64(pointsPerSide), float64(dimension))
	return Sum(pointsPerSide, dimension, fn) / N
}

// Find the minimum of fn over all Brillouin zone points.
func Min(pointsPerSide int, dimension int, fn BzFunc) float64 {
	minimum := func(next, min float64) float64 {
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
	points := bzPoints(L, d)
	total := start
	for i := 0; i < len(points); i++ {
		k := points[i]
		total = combine(fn(k), total)
	}
	return total
}

var pointsCache map[int]map[int][]vec.Vector = make(map[int]map[int][]vec.Vector)

// Produce a a slice of vectors whose values cover each first Brillouin zone
// point once. 
func bzPoints(L, d int) []vec.Vector {
	cachedL, okL := pointsCache[L]
	if okL {
		cachedD, okD := cachedL[d]
		if okD {
			return cachedD
		}
	} else {
		pointsCache[L] = make(map[int][]vec.Vector)
	}
	points := make([]vec.Vector, pow(L, d))
	// start is the minumum value of any component of a point
	start := -math.Pi
	// (finish - step) is the maximum value of any component of a point
	finish := -start
	// step is the separation between point components
	step := (finish - start) / float64(L)
	k := vec.ZeroVector(d)
	kIndex := make([]int, d)
	// set initial value for k
	for i := 0; i < d; i++ {
		k[i] = start
		kIndex[i] = 0
	}
	// iterate over Brillouin zone
	done := false
	for i := 0; !done; i++ {
		points[i] = vec.ZeroVector(len(k))
		copy(points[i], k)
		done = bzAdvance(k, kIndex, start, step, L, d)
	}
	pointsCache[L][d] = points
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

// Return x^y
func pow(x, y int) int64 {
	r := int64(1)
	for i := 0; i < y; i++ {
		r *= int64(x)
	}
	return r
}
