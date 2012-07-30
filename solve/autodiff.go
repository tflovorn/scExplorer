package solve

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_errno.h>
#include <gsl/gsl_deriv.h>
extern double go_val(double, void*);

static int centralDeriv(void *uservar, double x, double h, double *result, double *abserr) {
	gsl_function F;
	F.function = &go_val;
	F.params = uservar;
	int err = gsl_deriv_central(&F, x, h, result, abserr);
	return err;
}
*/
import "C"
import (
	"fmt"
	"math"
	"unsafe"
)
import vec "../vector"

type fnWithIndex struct {
	fn vec.FnDim0
	v  vec.Vector
	i  int
}

// Numerical central derivative of fn(v) with respect to v_i within tolerance
// epsabs. h is the initial step size.
func Derivative(fn vec.FnDim0, v vec.Vector, i int, h, epsabs float64) (float64, error) {
	iters, maxIters := 0, 100
	// results are bad for h too small or too large; we iterate in both directions
	hMin := h / 1e6
	hMax := h * 1e6
	hRising := true
	hInitial := h
	hOk := func(h float64) bool {
		if !hRising && h > hMin {
			return true
		} else if hRising && h < hMax {
			return true
		}
		return false
	}
	hAdvance := func(h float64) float64 {
		if hRising {
			if h*2.0 > hMax {
				hRising = false
				return hInitial / 2.0
			}
			return h * 2.0
		}
		return h / 2.0
	}
	fwi := unsafe.Pointer(&fnWithIndex{fn, v, i})
	x := C.double(v[i])
	result, abserr := C.double(0.0), C.double(math.MaxFloat64)
	for float64(abserr) > epsabs && iters < maxIters && hOk(h) {
		err := C.centralDeriv(fwi, x, C.double(h), &result, &abserr)
		if err != C.GSL_SUCCESS {
			err_str := C.GoString(C.gsl_strerror(err))
			return float64(result), fmt.Errorf("error in Derivative (GSL): %v\n", err_str)
		}
		iters++
		h = hAdvance(h)
	}
	return float64(result), nil
}

// Evaluate the go function contained in uservar at x.
//export go_val
func go_val(x C.double, uservar unsafe.Pointer) C.double {
	f := (*fnWithIndex)(uservar)
	f.v[f.i] = float64(x)
	val, err := f.fn(f.v)
	if err != nil {
		// TODO: handle error
	}
	return C.double(val)
}
