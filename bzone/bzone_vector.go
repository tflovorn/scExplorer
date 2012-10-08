package bzone

import "math"
import vec "../vector"

type BzVectorFunc func(k vec.Vector, out *vec.Vector)
type bzVectorConsumer func(next vec.Vector, total *vec.Vector)

func VectorAvg(pointsPerSide, gridDim, fnDim int, fn BzVectorFunc) vec.Vector {
	N := math.Pow(float64(pointsPerSide), float64(gridDim))
	sum := VectorSum(pointsPerSide, gridDim, fnDim, fn)
	avg := vec.ZeroVector(fnDim)
	for i := 0; i < fnDim; i++ {
		avg[i] = sum[i] / N
	}
	return avg
}

func VectorSum(pointsPerSide, gridDim, fnDim int, fn BzVectorFunc) vec.Vector {
	c := vec.ZeroVector(fnDim)
	add := func(next vec.Vector, total *vec.Vector) {
		for i := 0; i < fnDim; i++ {
			x := (*total)[i]
			y := next[i] - c[i]
			t := x + y
			c[i] = (t - x) - y
			(*total)[i] = t
		}
	}
	start := vec.ZeroVector(fnDim)
	return bzVectorReduce(add, start, pointsPerSide, gridDim, fnDim, fn)
}

func bzVectorReduce(combine bzVectorConsumer, start vec.Vector, L, d, fnDim int, fn BzVectorFunc) vec.Vector {
	points := bzPoints(L, d)
	total := start
	out := vec.ZeroVector(fnDim)
	for i := 0; i < len(points); i++ {
		k := points[i]
		fn(k, &out)
		combine(out, &total)
	}
	return total
}
