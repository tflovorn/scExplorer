package tempPair

import "github.com/tflovorn/scExplorer/tempAll"

func PairTempEnvironment(jsonData string) (*tempAll.Environment, error) {
	env, err := tempAll.NewEnvironment(jsonData)
	if err != nil {
		return nil, err
	}
	env.F0 = 0.0
	env.Alpha = -1

	return env, nil
}
