package plots

import (
	"flag"
	"os"
	"testing"
)

var testPlot = flag.Bool("testPlot", false, "Run test involving plots")

func TestPlotMPLParabola(t *testing.T) {
	flag.Parse()
	if !*testPlot {
		return
	}
	xs := []float64{1, 2, 3, 4, 5}
	ys := []float64{1, 4, 9, 16, 25}
	data := []Series{Series{xs, ys}}
	wd, _ := os.Getwd()
	params := map[string]string{FILE_KEY: wd + "/deleteme.mpljson_test_data", XLABEL_KEY: "$X$", YLABEL_KEY: "$Y$"}
	seriesParams := []map[string]string{map[string]string{}}
	err := PlotMPL(data, params, seriesParams, wd+"/grapher.py")
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
}
