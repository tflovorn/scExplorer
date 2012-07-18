package solve

/*
#cgo LDFLAGS: -lgsl
#include <gsl/gsl_vector.h>
#include <gsl/gsl_multiroots.h>
extern double go_f(const gsl_vector * x, void * fn);
extern void go_df(const gsl_vector * x, void * fn, gsl_vector * g);
extern void go_fdf(const gsl_vector * x, void * fn, double * f, gsl_vector * g);

*/
import "C"

import vec "../vector"

// Multidimensional root-finder. Implemented by providing an interface to GSL
// implementation of Powell's Hybrid method (gsl_multiroot_fdfsolver_hybridsj).
// Callback passing through cgo follows the model at:
// http://stackoverflow.com/questions/6125683/call-go-functions-from-c/6147097#6147097
func MultiDim(fn Diffable, start vec.Vector) (vec.Vector, error) {
	return []float64{}, nil
}
