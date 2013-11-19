package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../integrate"
	"../tempAll"
)

// Calculate the integral of y from 0 to ymax of: F(y) / (4*pi^2*a*sqrt(b)).
// a, b are parameters in pair spectrum: omega_+(k) = a*(kx^2 + ky^2) + b*(kz^2)
func OmegaIntegralY(env *tempAll.Environment, a, b float64, F func(float64) float64) (float64, error) {
	if a == 0.0 || b == 0.0 {
		return 0.0, nil
	}
	ymax := env.Beta * (-2.0*env.Mu_h + env.Mu_b)
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

// Calculate the integral of y from 0 to ymax of: F(y) / (8*pi^2*a).
// a, b are parameters in pair spectrum: omega_+(k) = a*(kx^2 + ky^2) + 2b*(1 - cos(kz))
func OmegaIntegralCos(env *tempAll.Environment, a, b float64, F func(float64, float64) float64) (float64, error) {
	if a == 0.0 || b == 0.0 {
		return 0.0, nil
	}
	// define the inner integral
	t := 1e-7
	integral_inner := func(kz float64) float64 {
		bterm := 2.0 * b * (1.0 - math.Cos(kz))
		ymax := env.Beta * (-2.0*env.Mu_h + env.Mu_b - bterm)
		if ymax <= 0.0 {
			return 0.0
		}
		if ymax/env.Beta > a*math.Pow(math.Pi, 2.0) {
			//fmt.Printf("ymax = %f is too large, replacing with (beta*a*pi^2)\n", ymax)
			//ymax = env.Beta * a * math.Pow(math.Pi, 2.0)
			panic(fmt.Errorf("ymax = %f is too large (a = %f, b = %f, env = %s), bailing from integral", ymax, a, b, env.String()))
		}
		innerF := func(y float64) float64 {
			return F(y, kz)
		}
		integral, abserr, err := integrate.Qags(innerF, 1e-10, ymax, t, t)
		if err != nil {
			fmt.Printf("inner integral error\n")
			panic(err)
		}
		if math.Abs(abserr) > t*100 {
			err = fmt.Errorf("inner integral too innaccurate (abserr = %e)", abserr)
			panic(err)
		}
		return integral
	}
	// calculate the full integral
	integral, abserr, err := integrate.Qags(integral_inner, -math.Pi, math.Pi, t, t)
	if err != nil {
		fmt.Printf("outer integral error\n")
		return 0.0, err
	}
	if math.Abs(abserr) > t*1000 {
		err = fmt.Errorf("outer integral too innaccurate (abserr = %e)", abserr)
		return 0.0, err
	}
	val := integral / (8.0 * math.Pow(math.Pi, 2.0) * env.Beta * a)
	return val, nil
}
