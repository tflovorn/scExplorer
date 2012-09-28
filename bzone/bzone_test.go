package bzone

import vec "../vector"

import (
	"fmt"
	"math"
	"testing"
)

// Check if the total number of points in the lattice is correct.
func TestTotalNumberOfPoints(t *testing.T) {
	// Return true if the number of points is correct on a lattice with 
	// side length L and dimension d.
	checkNumPoints := func(L, d int) (int, bool) {
		expected := pow(L, d)
		points := bzPoints(L, d)
		return len(points), len(points)-int(expected) == 0
	}
	N := 4
	count, correct := checkNumPoints(N, 2)
	if !correct {
		msg := fmt.Sprintf("Incorrect number of points from bzPoints(%d,%d)\n", N, 2)
		msg += fmt.Sprintf("Got %d, expected %d\n", count, int64(N*N))
		t.Fatal(msg)
	}
}

// Sum over product of sin(k_i) terms should be 0.
func TestSinSum(t *testing.T) {
	// Check if sum_k sum_i=0^d sin(k_i) = 0 within tolerance epsilon.
	checkSinSum := func(L, d int, epsilon float64) (float64, bool) {
		fn := func(k vec.Vector) float64 {
			r := 1.0
			for i := 0; i < d; i++ {
				r *= math.Sin(k[i])
			}
			return r
		}
		val := Sum(L, d, fn)
		return val, math.Abs(val) < epsilon
	}
	val, ok := checkSinSum(64, 2, 1e-9)
	if !ok {
		t.Fatalf("Incorrect sin sum value; got %f, expected 0.0", val)
	}
}

// Equivalent to checking number of points.
func TestOneSum(t *testing.T) {
	checkOneSum := func(L, d int, epsilon float64) (float64, bool) {
		expected := pow(L, d)
		fn := func(k vec.Vector) float64 {
			return 1.0
		}
		val := Sum(L, d, fn)
		return val, math.Abs(val-float64(expected)) < epsilon
	}
	val, ok := checkOneSum(4, 2, 1e-9)
	if !ok {
		t.Fatalf("Incorrect one sum value; got %f, expected %f", val, 64*64)
	}
}

// Does average return the correct average?
func TestOneAvg(t *testing.T) {
	checkOneSum := func(L, d int, epsilon float64) (float64, bool) {
		expected := 1.0
		fn := func(k vec.Vector) float64 {
			return 1.0
		}
		val := Avg(L, d, fn)
		return val, math.Abs(val-float64(expected)) < epsilon
	}
	val, ok := checkOneSum(4, 2, 1e-9)
	if !ok {
		t.Fatalf("Incorrect one average value; got %f, expected %f", val, 1.0)
	}

}

// Does Min return the correct minimum?
func TestMinQuadratic(t *testing.T) {
	expected := 0.0
	checkMin := func(L, d int, epsilon float64) (float64, bool) {
		fn := func(k vec.Vector) float64 {
			return k[0]*k[0] + k[1]*k[1]
		}
		val := Min(L, d, fn)
		return val, math.Abs(val-float64(expected)) < epsilon
	}
	val, ok := checkMin(16, 2, 1e-9)
	if !ok {
		t.Fatalf("Min reported incorrect minimum (got %v, expected %v)", val, expected)
	}
}
