package bzone

import "math"
import vec "../vector"

type BzVectorFunc func(k vec.Vector) vec.Vector
type bzVectorConsumer func(next, total vec.Vector) vec.Vector

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
	add := func(next, total vec.Vector) vec.Vector {
		r := vec.ZeroVector(fnDim)
		for i := 0; i < fnDim; i++ {
			y := next[i] - c[i]
			t := total[i] + y
			c[i] = (t - total[i]) - y
			r[i] = t
		}
		return r
	}
	start := vec.ZeroVector(fnDim)
	return bzVectorReduce(add, start, pointsPerSide, gridDim, fn)
}

func bzVectorReduce(combine bzVectorConsumer, start vec.Vector, L, d int, fn BzVectorFunc) vec.Vector {
	points := bzPoints(L, d)
	total := start
	for i := 0; i < len(points); i++ {
		k := points[i]
		total = combine(fn(k), total)
	}
	return total
}
