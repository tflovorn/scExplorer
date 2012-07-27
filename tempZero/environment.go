package tempZero

import "../tempAll"

type EnvZero struct {
	*tempAll.EnvAll

	F0    float64 // superconducting order parameter
	Alpha int     // SC gap symmetry parameter (s-wave = +1, d-wave = -1)
}
