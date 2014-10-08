package solve

import (
	"math"
	"testing"
)
import vec "github.com/tflovorn/scExplorer/vector"

func TestIterativeRosenbrock(t *testing.T) {
	epsAbs := []float64{1e-9, 1e-9}
	epsRel := []float64{1e-9, 1e-9}
	expected := []vec.Vector{[]float64{1.0}, []float64{1.0}}
	start := []vec.Vector{[]float64{10.0}, []float64{5.0}}
	stages, accept := rosenbrockSystems(1.0, 10.0, start)
	solutions, err := Iterative(stages, start, epsAbs, epsRel, accept)
	if err != nil {
		t.Fatal(err)
	}
	for i, vec := range expected {
		for j, val := range vec {
			if math.Abs(solutions[i][j]-val) > epsAbs[i] {
				t.Fatalf("incorrect Rosenbrock(1.0, 1.0) solution: %v", solutions)
			}
		}
	}
}

func rosenbrockSystems(a, b float64, start []vec.Vector) ([]DiffSystem, func([]vec.Vector)) {
	var x1, x2 float64
	accept := func(x []vec.Vector) {
		x1 = x[0][0]
		x2 = x[1][0]
	}
	accept(start)
	f1 := func(v vec.Vector) (float64, error) {
		return a * (1.0 - v[0]), nil
	}
	f2 := func(v vec.Vector) (float64, error) {
		return b * (v[0] - x1*x1), nil
	}
	d1 := SimpleDiffable(f1, 1, 1e-4, 1e-9)
	d2 := SimpleDiffable(f2, 1, 1e-4, 1e-9)
	s1 := Combine([]Diffable{d1})
	s2 := Combine([]Diffable{d2})
	return []DiffSystem{s1, s2}, accept
}
