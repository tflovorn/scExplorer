package fit

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <gsl/gsl_multifit.h>
*/
import "C"
import vec "../vector"

// Return the value for vector c which best fits the linear model y = Xc.
// X is a matrix given by its row vectors: each row corresponds to a y value;
// columns within the row correspond to coefficients for parameters in c.
func Linear(y vec.Vector, X []vec.Vector) vec.Vector {
	n := C.size_t(len(y))    // number of observations
	p := C.size_t(len(X[0])) // number of parameters
	c, gslY := C.gsl_vector_alloc(p), C.gsl_vector_alloc(n)
	cov, gslX := C.gsl_matrix_alloc(p, p), C.gsl_matrix_alloc(n, p)
	vecToGSL(y, gslY)
	matrixToGSL(X, gslX)
	chisq := C.double(0)
	work := C.gsl_multifit_linear_alloc(n, p)
	C.gsl_multifit_linear(gslX, gslY, c, cov, &chisq, work)
	C.gsl_multifit_linear_free(work)
	result := vecFromGSL(c)
	C.gsl_matrix_free(gslX)
	C.gsl_matrix_free(cov)
	C.gsl_vector_free(gslY)
	C.gsl_vector_free(c)
	return result
}

// Convert m to GSL format.
func matrixToGSL(m []vec.Vector, target *C.gsl_matrix) {
	numRows := len(m)
	numCols := len(m[0])
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			it, jt := C.size_t(i), C.size_t(j)
			C.gsl_matrix_set(target, it, jt, C.double(m[i][j]))
		}
	}
}
