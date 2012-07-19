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

func TestAdd(t *testing.T) {
	dim := 3
	v := ZeroVector(dim)
	v = v.Add([]float64{5.0, 4.0, 3.0})
	if v[0] != 5.0 || v[1] != 4.0 || v[2] != 3.0 {
		t.Fatalf("TestAdd added incorrect values; v = %v", v)
	}
}
