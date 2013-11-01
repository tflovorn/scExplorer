package integrate

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_errno.h>
#include <gsl/gsl_integration.h>
extern double integrate_go_val(double, void*);

static int gslQags(void *uservar, double a, double b, double epsabs, double epsrel, double *result, double *abserr) {
	int limit = 10000;
	gsl_integration_workspace *w = gsl_integration_workspace_alloc(limit);
	gsl_function F;
	F.function = &integrate_go_val;
	F.params = uservar;
	int err = gsl_integration_qags(&F, a, b, epsabs, epsrel, limit, w, result, abserr);
	gsl_integration_workspace_free(w);
	return err;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type floatFunc func(float64) float64

// Interface to gsl_integration_qags: "adaptive integration with [integrable]
// singularities". Integrate f from a to b and return the value of the
// integral and the absolute error.
func Qags(fn func(float64) float64, a, b, epsabs, epsrel float64) (result, absErr float64, err error) {
	cfn := unsafe.Pointer(&fn)
	C_result, C_absErr := C.double(0.0), C.double(0.0)
	// guard against panics during integration
	defer func() {
		if x := recover(); x != nil {
			result = 0.0
			absErr = 0.0
			err = x.(error)
		}
	}()
	// perform integration
	C_err := C.gslQags(cfn, C.double(a), C.double(b), C.double(epsabs), C.double(epsrel), &C_result, &C_absErr)
	if C_err != C.GSL_SUCCESS {
		err_str := C.GoString(C.gsl_strerror(C_err))
		return 0.0, 0.0, fmt.Errorf("error in Qags (GSL): %v\n", err_str)
	}
	result = float64(C_result)
	absErr = float64(C_absErr)
	return
}

//export integrate_go_val
func integrate_go_val(x C.double, uservar unsafe.Pointer) C.double {
	f := (*floatFunc)(uservar)
	val := (*f)(float64(x))
	return C.double(val)
}
