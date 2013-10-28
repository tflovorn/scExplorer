package integrate

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <stdlib.h>
#include <gsl/gsl_spline.h>
extern void integrate_spline_go_slice(int, void *, double *);

static double integrate_spline(void * xs, void * ys, int N, double a, double b) {
	gsl_interp_accel *acc = gsl_interp_accel_alloc();
	gsl_spline *spline = gsl_spline_alloc(gsl_interp_cspline, N);
	// convert xs and ys to C arrays
	double * cxs = malloc(N * sizeof(double));
	double * cys = malloc(N * sizeof(double));
	int i;
	for (i = 0; i < N; i++) {
		double vx = 0.0;
		double vy = 0.0;
		integrate_spline_go_slice(i, xs, &vx);
		integrate_spline_go_slice(i, ys, &vy);
		cxs[i] = vx;
		cys[i] = vy;
	}
	// do integration
	gsl_spline_init(spline, cxs, cys, N);
	double result = gsl_spline_eval_integ(spline, a, b, acc);
	// clean up
	gsl_spline_free(spline);
	gsl_interp_accel_free(acc);
	free(cxs);
	free(cys);
	return result;
}
*/
import "C"
import (
	"fmt"
	"sort"
	"unsafe"
)

// Integrate a function from a to b using cubic spline interpolation to
// approximate the function specified at the points xs with values ys.
// xs must be sorted and of the same length as ys, and (a, b) must be
// within (xs[0], xs[len(xs)-1]).
func Spline(xs, ys []float64, a, b float64) (float64, error) {
	// check preconditions
	N := len(xs)
	if N != len(ys) {
		return 0.0, fmt.Errorf("xs and ys must have the same length for spline interpolation")
	}
	if !sort.Float64sAreSorted(xs) {
		return 0.0, fmt.Errorf("xs values must be sorted for spline interpolation")
	}
	if a < xs[0] || b > xs[len(xs)-1] {
		return 0.0, fmt.Errorf("Spline integration bounds must be within xs")
	}
	// do spline integration
	val := C.integrate_spline(unsafe.Pointer(&xs), unsafe.Pointer(&ys), C.int(N), C.double(a), C.double(b))
	return float64(val), nil
}

//export integrate_spline_go_slice
func integrate_spline_go_slice(i C.int, xs unsafe.Pointer, result *C.double) {
	gxs := *((*[]float64)(xs))
	*result = C.double(gxs[i])
}
