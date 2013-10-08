package bessel

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <gsl/gsl_sf_bessel.h>
*/
import "C"

func ModifiedBesselFirstKindZeroth(x float64) float64 {
	return float64(C.gsl_sf_bessel_I0(C.double(x)))
}
