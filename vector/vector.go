package vector

type Vector []float64
type VectorFunc func(v Vector) float64

func ZeroVector(dimension int) Vector {
	v := make([]float64, dimension)
	return v
}

func (v Vector) Add(u Vector) Vector {
	if len(v) != len(u) {
		panic("cannot add vectors of different lengths")
	}
	r := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		r[i] = v[i] + u[i]
	}
	return r
}

func (v Vector) Mul(x float64) Vector {
	xv := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		xv[i] = x * (v[i])
	}
	return xv
}
