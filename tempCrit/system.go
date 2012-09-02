package tempCrit

import (
	"../solve"
	"../tempAll"
	"../tempPair"
	vec "../vector"
)

func CritTempStages(env *tempAll.Environment) ([]solve.DiffSystem, []vec.Vector, func([]vec.Vector)) {
	vars0 := []string{"D1", "Mu_h"}
	vars1 := []string{"Beta"}
	diffD1 := tempPair.AbsErrorD1(env, vars0)
	diffMu_h := tempPair.AbsErrorBeta(env, vars0)
	system0 := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	diffBeta := AbsErrorBeta(env, vars1)
	system1 := solve.Combine([]solve.Diffable{diffBeta})
	stages := []solve.DiffSystem{system0, system1}
	start := []vec.Vector{[]float64{env.D1, env.Mu_h}, []float64{env.Beta}}
	accept := func(x []vec.Vector) {
		env.D1 = x[0][0]
		env.Mu_h = x[0][1]
		env.Beta = x[1][0]
	}
	return stages, start, accept
}
