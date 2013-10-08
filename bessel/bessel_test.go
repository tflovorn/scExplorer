package bessel

import (
	"testing"
)

func TestModifiedBesselFirstKindZeroth(t *testing.T) {
	if ModifiedBesselFirstKindZeroth(0.0) != 1.0 {
		t.Error("incorrect Bessel function value for I0(0.0)")
	}
	if ModifiedBesselFirstKindZeroth(2.0) != 2.2795853023360673 {
		t.Error("incorrect Bessel function value for I0(2.0)")
	}
	if ModifiedBesselFirstKindZeroth(3.0) != 4.880792585865024 {
		t.Error("incorrect Bessel function value for I0(3.0)")
	}
}
