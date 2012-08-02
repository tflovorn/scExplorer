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
	} else if len(varNames) == 2 {
		return extractXY(dataSet, varNames[0], varNames[1])
	}
	// iterate through dataSet to create a map x->y for each z
	maps := make(map[float64]map[float64]float64)
	zs := make([]float64, 0)
	for _, data := range dataSet {
		val := reflect.ValueOf(data)
		x := val.FieldByName(varNames[0]).Float()
		y := val.FieldByName(varNames[1]).Float()
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
		for x, _ := range zmap {
			xs = append(xs, x)
		}
		sort.Float64s(xs)
		ys := make([]float64, len(xs))
		for j, x := range xs {
			ys[j] = zmap[x]
		}
		ret[i] = Series{xs, ys}
	}
	return ret
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
