package tempLow

import "../tempAll"

func Environment(jsonData string) (*tempAll.Environment, error) {
	env, err := tempAll.NewEnvironment(jsonData)
	if err != nil {
		return nil, err
	}
	env.Alpha = -1
	env.Mu_b = 0.0

	return env, nil
}
