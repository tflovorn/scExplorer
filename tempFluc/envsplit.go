package tempFluc

import (
	"github.com/tflovorn/scExplorer/tempAll"
	"github.com/tflovorn/scExplorer/tempCrit"
)

func EnvSplitTcB(baseEnv *tempAll.Environment, TcFactors, BeFields []float64, epsAbs, epsRel float64) ([]*tempAll.Environment, error) {
	TcEnv := baseEnv.Copy()
	TcEnv.Be_field = 0.0
	TcEnv.Mu_b = 0.0
	_, err := tempCrit.CritTempSolve(TcEnv, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	Tc := 1.0/TcEnv.Beta
	omegaFit, err := tempCrit.OmegaFit(TcEnv, tempCrit.OmegaPlus)
	if err != nil {
		return nil, err
	}
	TcEnv.A, TcEnv.B = omegaFit[0], omegaFit[2]
	TcEnv.PairCoeffsReady = true

	result := []*tempAll.Environment{}
	for _, TcFactor := range TcFactors {
		env := TcEnv.Copy()
		T := TcFactor * Tc
		env.Beta = 1.0/T
		env.Temp = T
		BeNum := len(BeFields)
		thisEnv_BeSplit := env.MultiSplit([]string{"Be_field"}, []int{BeNum}, []float64{BeFields[0]}, []float64{BeFields[BeNum - 1]})
		result = append(result, thisEnv_BeSplit...)
	}
	return result, nil
}
