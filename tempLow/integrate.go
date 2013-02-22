package tempLow

import (
	"fmt"
	"math"
)
import (
	"../integrate"
	"../tempAll"
)

// Calculate the integral of y from 0 to ymax of: F(y) / (4*pi^2*a*sqrt(b)).
func OmegaIntegralY(env *tempAll.Environment, F func(float64) float64) (float64, error) {
	a, b := env.A, env.B
	if a == 0.0 || b == 0.0 {
		return 0.0, nil
	}
	ymax := env.Beta * (-2.0 * env.Mu_h)
	if ymax <= 0.0 {
		return 0.0, nil
	}
	upper_a := math.Sqrt(ymax / (2.0 * env.Beta * a))
	if upper_a > math.Pi {
		ymax = 2.0 * env.Beta * a * math.Pow(math.Pi, 2.0)
	}
	upper_b := math.Sqrt(ymax / (env.Beta * b))
	if upper_b > math.Pi {
		ymax = env.Beta * b * math.Pow(math.Pi, 2.0)
	}
	t := 1e-7
	integral, abserr, err := integrate.Qags(F, 0.0, ymax, t, t)
	if err != nil {
		fmt.Printf("err conditions: ymax = %f; upper_a = %f; upper_b = %f\n", ymax, upper_a, upper_b)
		return 0.0, err
	}
	if math.Abs(abserr) > t*10 {
		err = fmt.Errorf("nu integral too innaccurate (abserr = %e)", abserr)
	}
	val := integral / (4.0 * math.Pow(math.Pi, 2.0) * a * math.Sqrt(b))
	return val, err
}
