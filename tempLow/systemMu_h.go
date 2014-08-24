package tempLow

import "fmt"
import (
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		//fmt.Printf("v=%v, variables=%v\n", v, variables)
		env.Set(v, variables)
		//fmt.Printf("env.Mu_h = %v\n", env.Mu_h)
		zv := vec.ZeroVector(3)
		omega0, err := Omega_pp(env, zv)
		if err != nil {
			fmt.Println(err)
			return 0.0, err
		}
		fmt.Printf("for v=%v got omega0=%f\n", v, omega0)
		return omega0, nil
	}
	h := 1e-6
	epsabs := 1e-6
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
