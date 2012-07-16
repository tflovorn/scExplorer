// One-dimensional zero-finding.
package solver

type OneDimFunc struct {
	f       func(x float64) float64
	epsilon float64 // maximum acceptable value for zero-finding
}

// min and max must bracket a zero.
func OneDimZero(fn OneDimFunc, min, max float64) (float64, error) {
	return 0.0, nil
}
