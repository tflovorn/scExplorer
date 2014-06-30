package tempLow

import "math"
import (
	"../tempCrit"
	"../tempAll"
	vec "../vector"
)

// Returns a function of omega which evaluates lambda_{r, s}(k, omega),
// where r and s are either +1 or -1.
// When omega = omega_{r, s}(k), lambda_{r, s}(k, omega) = 0.
// Imaginary part of M^D(p) should be 0 since we let
// i*omega -> omega + i*eps and set eps = 0.
//
// lambda_{r, s} = ((Re M^D_r)^2 - (M^{D,A}_r)^2)^(1/2) + i*s*Im(M^D_r)
func LambdaFn(env *tempAll.Environment, k vec.Vector, r, s int) func(float64) (float64, error) {
	return func(omega float64) (float64, error) {
		u, v := parts_MDiag(env, k, omega)
		Re_MD := u + float64(r)*v
		u_anom, v_anom := parts_MDiagAnom(env, k, omega)
		MDA := u_anom + float64(r)*v_anom
		return math.Pow(math.Pow(Re_MD, 2.0) - math.Pow(MDA, 2.0), 0.5), nil
	}
}

func parts_MDiag(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	cx, cy, cz := math.Cos(k[0]), math.Cos(k[1]), math.Cos(k[2])
	Ex := 2.0 * (env.T0*cy + env.Tz*cz)
	Ey := 2.0 * (env.T0*cx + env.Tz*cz)
	Pis := tempCrit.Pi(env, []float64{k[0], k[1]}, omega) // Pi_{xx, xy, yy}
	u := -0.25 * (-(1.0/Ex + 1.0/Ey) + Pis[0] + Pis[2])
	v := -0.5 * Pis[1]
	return u, v
}

func parts_MDiagAnom(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	PiAs_plus := PiAnom(env, []float64{k[0], k[1]}, omega) // Pi^A_{xx, xy, yy}
	PiAs_minus := PiAnom(env, []float64{-k[0], -k[1]}, -omega)
	u := -0.25 * (PiAs_plus[0] + PiAs_minus[0] + PiAs_plus[2] + PiAs_minus[2])
	v := -0.5 * (PiAs_plus[1] + PiAs_minus[1])
	return u, v
}
