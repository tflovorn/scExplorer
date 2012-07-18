package vector

import "testing"

func TestZeroVector(t *testing.T) {
	dim := 5
	v := ZeroVector(dim)
	for i := 0; i < dim; i++ {
		if v[i] != 0.0 {
			t.Fatalf("%d'ith component of zero vector is nonzero (value %v)", i, v[i])
		}
	}
}
