package tempLow

//import "math"
import "fmt"
import (
	//	"../bzone"
	"../solve"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		plusCoeffs, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
		if err != nil {
			return 0.0, err
		}
		fmt.Printf("mu_h plusCoeffs: %v\n", plusCoeffs)
		mu_pair := plusCoeffs[3]
		return mu_pair, nil
	}
	h := 2e-5
	epsabs := 1e-2
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
