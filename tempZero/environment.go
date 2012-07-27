package tempZero

import "math"
import (
	"../tempAll"
)

// Use serialized data to create an Environment appropriate for T=0.
func ZeroTempEnvironment(jsonData string) (*tempAll.Environment, error) {
	// initialize env with input data
	env, err := tempAll.NewEnvironment(jsonData)
	if err != nil {
		return nil, err
	}
	// set zero-temperature requirements
	env.Mu_b = 0.0
	env.Beta = math.Inf(1)

	return env, nil
}
