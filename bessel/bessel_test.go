package bessel

import (
	"testing"
)

func TestModifiedBesselFirstKindZeroth(t *testing.T) {
	input := []float64{0.0, 2.0, 3.0}
	expected := []float64{1.0, 2.2795853023360673, 4.880792585865024}
	for i, val := range(input) {
		b := ModifiedBesselFirstKindZeroth(val)
		if b != expected[i] {
			t.Errorf("incorrect Bessel function value for I0(%f); got %v, expected %v", val, b, expected[i])
		}
	}
}

func TestModifiedBesselFirstKindFirst(t *testing.T) {
	input := []float64{0.0, 2.0, 3.0}
	expected := []float64{0.0, 1.590636854637329, 3.9533702174026093}
	for i, val := range(input) {
		b := ModifiedBesselFirstKindFirst(val)
		if b != expected[i] {
			t.Errorf("incorrect Bessel function value for I0(%f); got %v, expected %v", val, b, expected[i])
		}
	}
}
