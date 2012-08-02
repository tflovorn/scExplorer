package plots

import (
	"reflect"
	"sort"
)

// A series of data to be used in a 2D graph. Implements the plotinum XYer
// interface.
type Series struct {
	xs, ys []float64
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

// Extract (x, y) points from dataSet where the names of x and y are given by
// varNames[0] and [1]. varNames[2] ("z") optionally specifies a variable to
// use to split the data into multiple series.
func ExtractSeries(dataSet []interface{}, varNames []string) []Series {
	if len(varNames) < 2 {
		panic("not enough variable names")
	}
	// iterate through dataSet to create a map x->y for each z
	maps := make(map[float64]map[float64]float64)
	zs := make([]float64, 0)
	no_z_xs, no_z_ymap := make([]float64, 0), make(map[float64]float64)
	for _, data := range dataSet {
		val := reflect.ValueOf(data)
		x := val.FieldByName(varNames[0]).Float()
		y := val.FieldByName(varNames[1]).Float()
		if len(varNames) > 2 {
			z := val.FieldByName(varNames[2]).Float()
			zmap, ok := maps[z]
			if !ok {
				zs = append(zs, z)
				maps[z] = make(map[float64]float64)
				zmap = maps[z]
			}
			zmap[x] = y
		} else {
			no_z_xs = append(no_z_xs, x)
			no_z_ymap[x] = y
		}
	}
	// if no z is specified, sort the single series and return
	if len(varNames) == 2 {
		sort.Float64s(no_z_xs)
		ys := make([]float64, len(no_z_xs))
		for i, x := range no_z_xs {
			ys[i] = no_z_ymap[x]
		}
		ret := make([]Series, 1)
		ret[0] = Series{no_z_xs, ys}
		return ret
	}
	// create the slice of Series in ascending-z order
	sort.Float64s(zs)
	ret := make([]Series, len(zs))
	for i, z := range zs {
		zmap := maps[z]
		// x should be in ascending order within the series
		xs := make([]float64, 0)
		for x, _ := range zmap {
			xs = append(xs, x)
		}
		sort.Float64s(xs)
		ys := make([]float64, len(xs))
		for j, x := range xs {
			ys[j] = x
		}
		ret[i] = Series{xs, ys}
	}
	return ret
}
