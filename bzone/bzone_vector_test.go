package bzone

import (
	"math"
	"testing"
)
import vec "github.com/tflovorn/scExplorer/vector"

func TestBzVectorAvgOnes(t *testing.T) {
	fn := func(k vec.Vector, out *vec.Vector) {
		(*out)[0] = 1.0
		(*out)[1] = 2.0
	}
	N := 8
	avg := VectorAvg(N, 2, 2, fn)
	epsilon := 1e-9
	if math.Abs(avg[0]-1.0) > epsilon || math.Abs(avg[1]-2.0) > epsilon {
		t.Fatalf("VectorAvg returned incorrect result")
	}
}

func TestBzVectorAvgSin(t *testing.T) {
	fn := func(k vec.Vector, out *vec.Vector) {
		(*out)[0] = math.Sin(k[0])
		(*out)[1] = math.Sin(k[1])
	}
	N := 16
	avg := VectorAvg(N, 2, 2, fn)
	epsilon := 1e-9
	if math.Abs(avg[0]-0.0) > epsilon || math.Abs(avg[1]-0.0) > epsilon {
		t.Fatalf("VectorAvg returned incorrect result")
	}
}
