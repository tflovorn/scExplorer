/*
package tempLow

import "math"
import (
	"github.com/tflovorn/scExplorer/bzone"
	"github.com/tflovorn/scExplorer/solve"
	"github.com/tflovorn/scExplorer/tempAll"
	vec "github.com/tflovorn/scExplorer/vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorF0(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := 1.0 / (env.T0 + env.Tz)
		//	rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerF0)) / 2.0
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerF0))
		return lhs - rhs, nil
	}
	h := 1e-5
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerF0(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) - math.Sin(k[1])
	E := env.BogoEnergy(k)
	return sxy * sxy * math.Tanh(env.Beta*E/2.0) / E
}
*/

package tempLow

import (
	"errors"
	"fmt"
)
import (
	"github.com/tflovorn/scExplorer/solve"
	"github.com/tflovorn/scExplorer/tempAll"
	"github.com/tflovorn/scExplorer/tempCrit"
	vec "github.com/tflovorn/scExplorer/vector"
)

func AbsErrorF0(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		if v.ContainsNaN() {
			fmt.Printf("got NaN in AbsErrorF0 (v=%v)\n", v)
			return 0.0, errors.New("NaN in input")
		}
		env.Set(v, variables)
		if !env.FixedPairCoeffs || !env.PairCoeffsReady {
			// Before we evaluate error in F0, Mu_h and D1 should have
			// appropriate values.
			eps := 1e-9
			_, err := D1MuSolve(env, eps, eps)
			if err != nil {
				return 0.0, err
			}
		}
		// F0 equation error = x - x1 - x2
		x1 := X1(env)
		x2, err := tempCrit.X2(env)
		if err != nil {
			fmt.Printf("error from X2(): %v\n", err)
			return 0.0, err
			//x2 = 0.0
		}
		lhs := env.X
		rhs := x1 + x2
		return lhs - rhs, nil
	}
	h := 1e-5
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
