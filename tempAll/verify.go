package tempAll

import (
	"fmt"
	"math"
	"reflect"
)
import "../solve"

// A Systemer returns a system for solving `env` and a starting point.
type Systemer func(env *Environment) (solve.DiffSystem, []float64)

func VerifySolution(env *Environment, sv Solver, st Systemer, vars []string, epsAbs, epsRel float64, expected []float64) error {
	solution, err := sv(env, epsAbs, epsRel)
	if err != nil {
		return err
	}
	// sv should leave env in solved state
	for i := 0; i < len(solution); i++ {
		// extract vars[i] from env: expect a float
		field := reflect.ValueOf(env).Elem().FieldByName(vars[i])
		v := field.Float()

		if math.Abs(solution[i]-v) > epsAbs {
			return fmt.Errorf("Env fails to match solution; env = %v; solution = %v", env, solution)
		}
	}
	// the solution we got should give 0 error within tolerances
	system, _ := st(env)
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		return fmt.Errorf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > epsAbs {
			return fmt.Errorf("error in pair temp system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > epsAbs {
			return fmt.Errorf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}

	return nil
}
