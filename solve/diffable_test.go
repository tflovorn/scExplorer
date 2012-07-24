package solve

import (
	"testing"
)

import vec "../vector"

func TestCombineRosenbrock(t *testing.T) {
	checkRosenbrock := func(a, b float64) {
		f1, f2 := RosenbrockF(a, b)
		df1, df2 := RosenbrockDf(a, b)
		fdf1, fdf2 := RosenbrockFdf(a, b)
		diff1 := Diffable{f1, df1, fdf1, 2}
		diff2 := Diffable{f2, df2, fdf2, 2}
		system := Combine([]Diffable{diff1, diff2})
		v := []float64{1.0, 1.0}
		vf1, _ := f1(v)
		vf2, _ := f2(v)
		vdf1, _ := df1(v)
		vdf2, _ := df2(v)
		svf, _ := system.F(v)
		svdf, _ := system.Df(v)
		if vf1 != svf[0] || vf2 != svf[1] || !vdf1.Equals(svdf[0]) || !vdf2.Equals(svdf[1]) {
			t.Fatalf("Diffable combine failed; vf1 = %v, vf2 = %v, vdf1 = %v, vdf2 = %v, svf = %v, svdf = %v", vf1, vf2, vdf1, vdf2, svf, svdf)
		}
	}
	checkRosenbrock(1.0, 10.0)
}

func RosenbrockF(a, b float64) (vec.FnDim0, vec.FnDim0) {
	f1 := func(v vec.Vector) (float64, error) {
		return a * (1.0 - v[0]), nil
	}
	f2 := func(v vec.Vector) (float64, error) {
		return b * (v[1] - v[0]*v[0]), nil
	}
	return f1, f2
}

func RosenbrockDf(a, b float64) (vec.FnDim1, vec.FnDim1) {
	df1 := func(v vec.Vector) (vec.Vector, error) {
		df1_x := -a
		df1_y := 0.0
		return []float64{df1_x, df1_y}, nil
	}
	df2 := func(v vec.Vector) (vec.Vector, error) {
		df2_x := -2.0 * b * v[0]
		df2_y := b
		return []float64{df2_x, df2_y}, nil
	}
	return df1, df2
}

func RosenbrockFdf(a, b float64) (vec.FnDim0_1, vec.FnDim0_1) {
	f1, f2 := RosenbrockF(a, b)
	df1, df2 := RosenbrockDf(a, b)
	fdf1 := func(v vec.Vector) (float64, vec.Vector, error) {
		vf1, _ := f1(v)
		vdf1, _ := df1(v)
		return vf1, vdf1, nil
	}
	fdf2 := func(v vec.Vector) (float64, vec.Vector, error) {
		vf2, _ := f2(v)
		vdf2, _ := df2(v)
		return vf2, vdf2, nil

	}
	return fdf1, fdf2
}
