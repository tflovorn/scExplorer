// Multidimensional minimization based on the gradient descent method.
package minimize

import . "../vector"

type GradientFunc struct {
	f       VectorFunc
	grad    []VectorFunc
	epsilon float64 // maximum acceptable value for zero-finding
}

func GradientMin(gfs []GradientFunc, start Vector) (Vector, error) {
	return []float64{}, nil
}

func Gradient(gf GradientFunc, v Vector) Vector {
	dimension := len(gf.grad)
	xs := make([]float64, dimension)
	for i := 0; i < dimension; i++ {
		xs[i] = gf.grad[i](v)
	}
	return xs
}

func oneDimInDirection(vf VectorFunc, start, dir Vector) func(float64) float64 {
	return func(alpha float64) float64 {
		return vf(start.Add(dir.Mul(alpha)))
	}
}
