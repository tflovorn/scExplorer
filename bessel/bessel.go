package bessel

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <gsl/gsl_sf_bessel.h>
*/
import "C"

// Zeroth-order modified Bessel function of the first kind (I0)
func ModifiedBesselFirstKindZeroth(x float64) float64 {
	return float64(C.gsl_sf_bessel_I0(C.double(x)))
}

// First-order modified Bessel function of the first kind (I1)
func ModifiedBesselFirstKindFirst(x float64) float64 {
	return float64(C.gsl_sf_bessel_I1(C.double(x)))

}
