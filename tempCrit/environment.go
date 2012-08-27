package tempCrit

import "../tempAll"

func CritTempEnvironment(jsonData string) (*tempAll.Environment, error) {
	env, err := tempAll.NewEnvironment(jsonData)
	if err != nil {
		return nil, err
	}
	env.F0 = 0.0
	env.Alpha = -1
	env.Mu_b = 0.0

	return env, nil
}
