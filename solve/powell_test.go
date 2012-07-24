package solve

import (
	"math"
	"testing"
)

func TestSolveRosenbrock(t *testing.T) {
	epsAbs := 1e-9
	system := RosenbrockSystem(1.0, 10.0)
	start := []float64{10.0, 5.0}
	solution, err := MultiDim(system, start, epsAbs, 1e-9)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(solution[0]-1.0) > epsAbs || math.Abs(solution[1]-1.0) > epsAbs {
		t.Fatalf("failed to produce correct Rosenbrock(1.0, 10.0) solution; got %v", solution)
	}
}

func RosenbrockSystem(a, b float64) DiffSystem {
	f1, f2 := RosenbrockF(a, b)
	df1, df2 := RosenbrockDf(a, b)
	fdf1, fdf2 := RosenbrockFdf(a, b)
	diff1 := Diffable{f1, df1, fdf1, 2}
	diff2 := Diffable{f2, df2, fdf2, 2}
	return Combine([]Diffable{diff1, diff2})
}
