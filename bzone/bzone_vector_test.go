package bzone

import (
	"math"
	"testing"
)
import vec "../vector"

func BzVectorAvgOnes(t *testing.T) {
	fn := func(k vec.Vector) vec.Vector {
		return []float64{1.0, 2.0}
	}
	N := 8
	avg := VectorAvg(N, 2, 2, fn)
	epsilon := 1e-9
	if math.Abs(avg[0]-1.0) > epsilon || math.Abs(avg[1]-2.0) > epsilon {
		t.Fatalf("VectorAvg returned incorrect result")
	}
}

func BzVectorAvgSin(t *testing.T) {
	fn := func(k vec.Vector) vec.Vector {
		return []float64{math.Sin(k[0]), math.Sin(k[1])}
	}
	N := 16
	avg := VectorAvg(N, 2, 2, fn)
	epsilon := 1e-9
	if math.Abs(avg[0]-0.0) > epsilon || math.Abs(avg[1]-0.0) > epsilon {
		t.Fatalf("VectorAvg returned incorrect result")
	}
}
