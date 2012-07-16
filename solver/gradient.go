// Multidimensional minimization based on the gradient descent method.
package solver

type Vector []float64
type VectorFunc func(v Vector) float64

type GradientFunc struct {
	f       VectorFunc
	grad    []VectorFunc
	epsilon float64 // maximum acceptable value for zero-finding
}

func GradientZero(gfs []GradientFunc, start Vector) (Vector, error) {
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
