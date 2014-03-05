package plots

import (
	"fmt"
	"reflect"
	"sort"
)

var COLOR_STYLES = []string{"k-", "r-", "b-", "g-", "c-", "m-", "y-", "k--", "r--", "b--", "g--", "c--", "m--", "y--"}
var PRINT_STYLES = []string{"k-", "k--", "k-.", "k:"}
var DEFAULT_STYLES = COLOR_STYLES

// Names of parameters relevant to MultiPlot
type GraphVars struct {
	// x and y axis of the plot
	X, Y string
	// additional parameters
	Params []string
	// Latex labels for additional parameters
	ParamLabels []string
	// if X/Y variable name is not given, use this to find value of X/Y
	XFunc func(interface{}) float64
	YFunc func(interface{}) float64
}

// Create a plot for each combination of vars.Params contained in data.
// graphParams[FILE_KEY] must specify a file path for output. grapherPath
// must specify the location of the Python graphing script.
func MultiPlot(data []interface{}, errs []error, vars GraphVars, graphParams map[string]string, grapherPath string) error {
	return multiPlotHelper(data, errs, vars, graphParams, grapherPath, DEFAULT_STYLES)
}

// Same as MultiPlot except that printStyle = true specifies to use line
// styles appropriate for printing in black & white.
func MultiPlotStyle(data []interface{}, errs []error, vars GraphVars, graphParams map[string]string, grapherPath string, printStyle bool) error {
	if printStyle {
		return multiPlotHelper(data, errs, vars, graphParams, grapherPath, PRINT_STYLES)
	}
	return multiPlotHelper(data, errs, vars, graphParams, grapherPath, COLOR_STYLES)
}

func multiPlotHelper(data []interface{}, errs []error, vars GraphVars, graphParams map[string]string, grapherPath string, seriesStyles []string) error {
	data = stripNils(data)
	// we need to know the values of vars.Params to set up constraint
	allParamValues, err := extractParamValues(data, errs, vars)
	if err != nil {
		return err
	}
	// iterate over combinations of parameters
	basePath := graphParams[FILE_KEY]
	primaryNames, primaryLabels, secondaries := paramCombinations(allParamValues, vars)
	for i := 0; i < len(primaryNames); i++ {
		// graphParams[FILE_KEY] needs to be modified for this combination of secondary params
		extraPath := ""
		for name, val := range secondaries[i] {
			extraPath = extraPath + fmt.Sprintf("%s_%f_", name, val)
		}
		graphParams[FILE_KEY] = basePath + extraPath

		series, primaryVals := ExtractSeries(data, errs, []string{vars.X, vars.Y, primaryNames[i]}, secondaries[i], vars.XFunc, vars.YFunc)
		sp := MakeSeriesParams(primaryLabels[i], "%.3f", primaryVals, seriesStyles)
		err := PlotMPL(series, graphParams, sp, grapherPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func stripNils(data []interface{}) []interface{} {
	ret := make([]interface{}, 0)
	for _, d := range data {
		if d != nil {
			ret = append(ret, d)
		}
	}
	return ret
}

func extractParamValues(data []interface{}, errs []error, vars GraphVars) (map[string][]float64, error) {
	// get values for vars.Params
	paramValues := make(map[string][]float64)
	for i, d := range data {
		if errs[i] != nil {
			continue
		}
		dv := reflect.ValueOf(d)
		for _, p := range vars.Params {
			pv := dv.FieldByName(p)
			if !pv.IsValid() {
				return paramValues, fmt.Errorf("Invalid parameter name %s", p)
			}
			pf := pv.Float()
			_, ok := paramValues[p]
			if !ok {
				paramValues[p] = make([]float64, 0)
			}
			if !contains(paramValues[p], pf) {
				paramValues[p] = append(paramValues[p], pf)
			}
		}
	}
	// sort param values
	for _, v := range paramValues {
		sort.Float64s(v)
	}

	return paramValues, nil
}

// Return true if xs contains val and false otherwise
func contains(xs []float64, val float64) bool {
	for _, x := range xs {
		if x == val {
			return true
		}
	}
	return false
}

// Return primaryNames, primaryLabels, secondaries. All variables get a
// chance to be primary; while a variable is primary all possible
// combinations of secondary variables are iterated through.
func paramCombinations(allParamValues map[string][]float64, vars GraphVars) ([]string, []string, []map[string]float64) {
	if len(vars.Params) == 1 {
		primaryNames := []string{vars.Params[0]}
		primaryLabels := []string{vars.ParamLabels[0]}
		secondaries := []map[string]float64{map[string]float64{}}
		return primaryNames, primaryLabels, secondaries
	}
	primaryNames := make([]string, 0)
	primaryLabels := make([]string, 0)
	secondaries := make([]map[string]float64, 0)
	for i, primaryName := range vars.Params {
		secondaryNames := without(vars.Params, primaryName)
		combos := combinations(allParamValues, secondaryNames)
		for _, secondary := range combos {
			primaryNames = append(primaryNames, primaryName)
			primaryLabels = append(primaryLabels, vars.ParamLabels[i])
			secondaries = append(secondaries, secondary)
		}
	}
	return primaryNames, primaryLabels, secondaries
}

// Return a copy of vals with drop removed
func without(vals []string, drop string) []string {
	ret := make([]string, 0)
	for _, v := range vals {
		if v != drop {
			ret = append(ret, v)
		}
	}
	return ret
}

// Return all possible combinations of named parameters given in vals
func combinations(vals map[string][]float64, names []string) []map[string]float64 {
	if len(names) == 0 {
		return nil
	}
	if len(names) == 1 {
		ret := make([]map[string]float64, 0)
		for _, val := range vals[names[0]] {
			thisMap := make(map[string]float64)
			thisMap[names[0]] = val
			ret = append(ret, thisMap)
		}
		return ret
	}
	remainingNames := without(names, names[0])
	remainingCombos := combinations(vals, remainingNames)
	ret := make([]map[string]float64, 0)
	for _, combo := range remainingCombos {
		for _, val := range vals[names[0]] {
			thisMap := copyMap(combo)
			thisMap[names[0]] = val
			ret = append(ret, thisMap)
		}
	}
	return ret
}

// Return a copy of the given map
func copyMap(m map[string]float64) map[string]float64 {
	ret := make(map[string]float64)
	for k, v := range m {
		ret[k] = v
	}
	return ret
}
