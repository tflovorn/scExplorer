package seriesaccel

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <stdlib.h>
#include <gsl/gsl_math.h>
#include <gsl/gsl_sum.h>
extern void levin_go_f(int, void *, double *);

static void levin(void * fn, int iStart, int Nterms, double * result, double * absErr, int * terms) {
	gsl_sum_levin_u_workspace * w = gsl_sum_levin_u_alloc(Nterms);
	double * t = malloc(Nterms * sizeof(double));
	int i;

	for (i = iStart; (i - iStart) < Nterms; i++) {
		double val = 0.0;
		levin_go_f(i, fn, &val);
		t[i-iStart] = val;
	}

	gsl_sum_levin_u_accel(t, Nterms, w, result, absErr);
	*terms = w->terms_used;

	free(t);
	gsl_sum_levin_u_free(w);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type TermFn func(int) float64

// Evaluate the series \sum_{i=iStart}^{\infty}[termFn(i)] using the
// Levin u-transform as provided by GSL. Returns the sum and error estimate.
// The initial evaluation will use Nstart terms. The series acceleration
// function uses as many as are necessary to minimize the error; if all terms
// are consumed, the number of terms used will be increased.
func Levin_u(fn TermFn, iStart, Nstart int) (float64, float64) {
	result, absErr := C.double(0.0), C.double(0.0)
	terms := C.int(0)

	C.levin(unsafe.Pointer(&fn), C.int(iStart), C.int(Nstart), &result, &absErr, &terms)
	// if number of terms used == Nstart, evaluate again with larger Nstart
	if int(terms) == Nstart {
		if Nstart > 300 {
			fmt.Printf("bailing early from Levin_u; result=%e, absErr=%e, terms=%d\n", float64(result), float64(absErr), Nstart)
			return float64(result), float64(absErr)
		}
		return Levin_u(fn, iStart, Nstart*2)
	}
	return float64(result), float64(absErr)
}

//export levin_go_f
func levin_go_f(i C.int, fn unsafe.Pointer, result *C.double) {
	gofn := *((*TermFn)(fn))
	*result = C.double(gofn(int(i)))
}
