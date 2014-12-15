package plots

import (
	"math"
	"reflect"
	"sort"
)

// A series of data to be used in a 2D graph. Implements the plotinum XYer
// interface.
type Series struct {
	xs, ys []float64
}

func MakeSeries(xs, ys []float64) Series {
	s := Series{xs, ys}
	return s
}

func (s *Series) X(i int) float64 {
	return s.xs[i]
}

func (s *Series) Y(i int) float64 {
	return s.ys[i]
}

func (s *Series) Len() int {
	return len(s.xs)
}

func (s *Series) Pairs() [][]float64 {
	L := s.Len()
	r := make([][]float64, L)
	for i := 0; i < L; i++ {
		r[i] = make([]float64, 2)
		r[i][0] = s.X(i)
		r[i][1] = s.Y(i)
	}
	return r
}

// Extract (x, y) points from dataSet where the names of x and y are given by
// varNames[0] and [1]. varNames[2] ("z") optionally specifies a variable to
// use to split the data into multiple series. `constraints` specifies
// parameter values to include; for example if constraints = {"Tz": 0.1}, only
// data points with Tz = 0.1 will be extracted. If the name of y is given as
// the empty string, `YFunc` is used to obtain a value for y instead.
func ExtractSeries(dataSet []interface{}, errs []error, varNames []string, constraints map[string]float64, XFunc, YFunc func(interface{}) float64, addZeros bool) ([]Series, []float64) {
	if len(varNames) < 2 {
		panic("not enough variable names for ExtractSeries")
	} else if len(varNames) == 2 {
		return extractXY(dataSet, varNames[0], varNames[1]), nil
	}
	// iterate through dataSet to create a map x->y for each z
	maps := make(map[float64]map[float64]float64)
	zs := make([]float64, 0)
	for i, data := range dataSet {
		if errs[i] != nil {
			continue
		}
		val := reflect.ValueOf(data)
		// check that constraints hold for this point
		ok := true
		for k, v := range constraints {
			c := val.FieldByName(k).Float()
			if v != c {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// constraints fit; keep point
		var x, y float64
		if varNames[0] != "" {
			x = val.FieldByName(varNames[0]).Float()
		} else {
			x = XFunc(data)
		}
		if varNames[1] != "" {
			yval := val.FieldByName(varNames[1])
			if yval.IsValid() {
				y = yval.Float()
			} else {
				y = 0.0
			}
		} else {
			y = YFunc(data)
		}
		if math.IsNaN(y) {
			// replace with default values
			// might want to be smarter about this
			y = 0.0
		}
		z := val.FieldByName(varNames[2]).Float()
		zmap, ok := maps[z]
		if !ok {
			zs = append(zs, z)
			maps[z] = make(map[float64]float64)
			zmap = maps[z]
		}
		zmap[x] = y
	}
	// create the slice of Series in ascending-z order
	sort.Float64s(zs)
	ret := make([]Series, len(zs))
	for i, z := range zs {
		zmap := maps[z]
		// x should be in ascending order within the series
		xs := make([]float64, 0)
		if addZeros {
			xs = append(xs, 0.0)
		}
		for x, _ := range zmap {
			xs = append(xs, x)
		}
		sort.Float64s(xs)
		ys := make([]float64, len(xs))
		if addZeros {
			ys[0] = 0.0
		}
		for j, x := range xs {
			if addZeros && j == 0 {
				continue
			}
			ys[j] = zmap[x]
		}
		ret[i] = Series{xs, ys}
	}
	return ret, zs
}

func extractXY(dataSet []interface{}, varX, varY string) []Series {
	xs, ymap := make([]float64, 0), make(map[float64]float64)
	for _, data := range dataSet {
		val := reflect.ValueOf(data)
		x := val.FieldByName(varX).Float()
		y := val.FieldByName(varY).Float()
		xs = append(xs, x)
		ymap[x] = y
	}
	sort.Float64s(xs)
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = ymap[x]
	}
	return []Series{Series{xs, ys}}
}
