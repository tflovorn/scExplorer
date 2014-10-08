package tempAll

import (
	"fmt"
	"math"
	"reflect"
)
import (
	"github.com/tflovorn/scExplorer/bzone"
	"github.com/tflovorn/scExplorer/serialize"
	vec "github.com/tflovorn/scExplorer/vector"
)

// Container for variables relevant at all temperatures.
type Environment struct {
	// Program parameters:
	PointsPerSide int // length of one side of the lattice

	// Constant physical parameters:
	X        float64 // average density of holons
	T0       float64 // nn one-holon hopping energy
	Thp      float64 // direct nnn one-holon hopping energy
	Tz       float64 // inter-planar one-holon hopping energy
	Alpha    int     // SC gap symmetry parameter (s-wave = +1, d-wave = -1)
	Be_field float64 // magnetic field (flux density) times e along the c axis (eB is unitless)

	// Dynamically determined physical parameters:
	D1   float64 // nnn hopping parameter generated by two-holon hopping
	Mu_h float64 // holon chemical potential

	// May be constant or dynamically determined:
	Beta float64 // inverse temperature
	F0   float64 // superconducting order parameter (0 if T >= Tc)
	Mu_b float64 // holon pair (bosonic) chemical potential (0 if T <= Tc)
	A, B float64 // pair spectrum parameters

	// Behavior flags:

	// Iterate solution for Mu_b in tempFluc.SolveD1Mu_hMu_b.
	IterateD1Mu_hMu_b bool
	// Use kz^2 in pair spectrum - incompatible with finite magnetic field.
	// If false, use cosine spectrum in tempCrit/tempFluc; cosine spectrum not implemented in tempLow.
	PairKzSquaredSpectrum bool
	// If true, include poles giving omega_- spectrum in the calculation.
	// Has no effect if PairKzSquaredSpectrum = false (TODO - extend to cos(kz) spectrum).
	// Expect negligible different between this being on or off.
	OmegaMinusPoles bool
	// Fix pair spectrum coefficients to their values at Tc (helps match F0 = 0 to Tc in T < Tc calculation).
	FixedPairCoeffs bool
	// If FixedPairCoeffs = true, stop varying pair spectrum coefficients after PairCoeffsReady is set to true.
	PairCoeffsReady bool

	// Cached values:
	epsilonMinCache  float64
	lastEpsilonMinD1 float64
}

type Wrappable func(*Environment, vec.Vector) float64

// ===== Utility functions =====

// Wrap fn with a function which depends only on a vector
func WrapFunc(env *Environment, fn Wrappable) bzone.BzFunc {
	return func(k vec.Vector) float64 {
		return fn(env, k)
	}
}

// Create an Environment from the given serialized data.
func NewEnvironment(jsonData string) (*Environment, error) {
	// initialize env with input data
	env := new(Environment)
	err := serialize.CopyFromJSON(jsonData, env)
	if err != nil {
		return nil, err
	}
	// hack to get around JSON's lack of support for Inf
	if env.Beta == math.MaxFloat64 {
		env.Beta = math.Inf(1)
	}
	// initialize cache
	env.setEpsilonMinCache()
	env.lastEpsilonMinD1 = env.D1

	return env, nil
}

// Convert to string by marshalling to JSON
func (env *Environment) String() string {
	if env.Beta == math.Inf(1) {
		// hack to get around JSON's choice to not allow Inf
		env.Beta = math.MaxFloat64
	}
	marshalled, err := serialize.MakeJSON(env)
	if err != nil {
		panic(err)
	}
	if env.Beta == math.MaxFloat64 {
		env.Beta = math.Inf(1)
	}
	return marshalled
}

// Create and return a copy of env
func (env *Environment) Copy() *Environment {
	marshalled := env.String()
	thisCopy, err := NewEnvironment(marshalled)
	if err != nil {
		// shouldn't get here (env should always be copyable)
		panic(err)
	}
	return thisCopy
}

// Iterate through v and vars simultaneously. vars specifies the names of
// fields to change in env (they are set to the values given in v).
// Panics if vars specifies a field not contained in env (or a field of
// non-float type).
func (env *Environment) Set(v vec.Vector, vars []string) {
	ev := reflect.ValueOf(env).Elem()
	for i := 0; i < len(vars); i++ {
		field := ev.FieldByName(vars[i])
		if field == reflect.Zero(reflect.TypeOf(env)) {
			panic(fmt.Sprintf("Field %v not present in Environment", vars[i]))
		}
		if field.Type().Kind() != reflect.Float64 {
			panic(fmt.Sprintf("Field %v is non-float", vars[i]))
		}
		field.SetFloat(v[i])
	}
}

// Split env into many copies with different values of the variable given by
// varName (N values running from min to max).
func (env *Environment) Split(varName string, N int, min, max float64) []*Environment {
	step := (max - min) / float64(N-1)
	if N == 1 {
		step = 0
	}
	rets := make([]*Environment, N)
	for i := 0; i < N; i++ {
		x := min + float64(i)*step
		thisCopy := env.Copy()
		thisCopy.Set([]float64{x}, []string{varName})
		rets[i] = thisCopy
	}
	return rets
}

// Call Split on each env in envs for each var in varNames to create a
// "Cartesian product" of the desired splits.
func (env *Environment) MultiSplit(varNames []string, Ns []int, mins, maxs []float64) []*Environment {
	if len(varNames) == 0 {
		return nil
	}
	oneSplit := env.Split(varNames[0], Ns[0], mins[0], maxs[0])
	if len(varNames) == 1 {
		return oneSplit
	}
	ret := make([]*Environment, 0)
	for _, osEnv := range oneSplit {
		ms := osEnv.MultiSplit(varNames[1:], Ns[1:], mins[1:], maxs[1:])
		for _, e := range ms {
			ret = append(ret, e)
		}
	}
	return ret
}

// ===== Physics functions =====

// Scaled hopping energy
func (env *Environment) Th() float64 {
	return env.T0 * (1.0 - env.X)
}

// Single-holon energy. Minimum is 0.
// env.EpsilonMin must be set to the value returned by EpsilonMin before
// calling this function.
func (env *Environment) Epsilon_h(k vec.Vector) float64 {
	return env.epsilonBar(k) - env.getEpsilonMin()
}

// Single-holon energy without fixed minimum.
func (env *Environment) epsilonBar(k vec.Vector) float64 {
	sx, sy := math.Sin(k[0]), math.Sin(k[1])
	return 2.0*env.Th()*((sx+sy)*(sx+sy)-1.0) + 4.0*(2.0*env.D1*env.T0-env.Thp)*sx*sy
}

// Get minimum value of env.Epsilon. If env.D1 hasn't changed since the last
// call to this function, return a cached value.
func (env *Environment) getEpsilonMin() float64 {
	if env.D1 != env.lastEpsilonMinD1 {
		env.setEpsilonMinCache()
		env.lastEpsilonMinD1 = env.D1
	}
	return env.epsilonMinCache
}

// Find the minimum of EpsilonBar.
func (env *Environment) setEpsilonMinCache() {
	worker := func(k vec.Vector) float64 {
		return env.epsilonBar(k)
	}
	env.epsilonMinCache = bzone.Min(env.PointsPerSide, 2, worker)
	//println(env.epsilonMinCache)
}

// Single-holon energy minus chemical potential. Minimum is -env.Mu_h.
func (env *Environment) Xi_h(k []float64) float64 {
	return env.Epsilon_h(k) - env.Mu_h
}

// Superconducting gap function.
func (env *Environment) Delta_h(k vec.Vector) float64 {
	return 4.0 * (env.T0 + env.Tz) * env.F0 * (math.Sin(k[0]) + float64(env.Alpha)*math.Sin(k[1]))
}

// Bogolyubov quasiparticle energy.
func (env *Environment) BogoEnergy(k vec.Vector) float64 {
	xi := env.Xi_h(k)
	delta := env.Delta_h(k)
	return math.Sqrt(xi*xi + delta*delta)
}

// Fermi distribution function.
func (env *Environment) Fermi(energy float64) float64 {
	if energy == 0.0 {
		return 0.5
	}
	// Temperature is 0 or e^(Beta*energy) is too big to calculate
	if env.Beta == math.Inf(1) || env.Beta >= math.Abs(math.MaxFloat64/energy) || math.Abs(env.Beta*energy) >= math.Log(math.MaxFloat64) {
		if energy <= 0 {
			return 1.0
		}
		return 0.0
	}
	// nonzero temperature
	return 1.0 / (math.Exp(energy*env.Beta) + 1.0)
}

// Extract the temperature from env
func GetTemp(data interface{}) float64 {
	env := data.(Environment)
	return 1.0 / env.Beta
}
