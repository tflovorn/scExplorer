package tempAll

import "strconv"
import "../plots"

func FermiSurface(env *Environment, outPrefix, grapherPath string) error {
	// blank data slots - relevant only for scatter plot
	data := []plots.Series{}
	seriesParams := []map[string]string{}
	// required data for FS plot
	params := make(map[string]string)
	params["plot_type"] = "Fermi_surface"
	params[plots.XLABEL_KEY] = "$k_x$"
	params[plots.YLABEL_KEY] = "$k_y$"
	params[plots.FILE_KEY] = outPrefix
	params["th"] = strconv.FormatFloat(env.Th(), 'f', 6, 64)
	params["thp"] = strconv.FormatFloat(env.Thp, 'f', 6, 64)
	params["t0"] = strconv.FormatFloat(env.T0, 'f', 6, 64)
	params["D1"] = strconv.FormatFloat(env.D1, 'f', 6, 64)
	params["Mu_h"] = strconv.FormatFloat(env.Mu_h, 'f', 6, 64)
	params["epsilon_min"] = strconv.FormatFloat(env.getEpsilonMin(), 'f', 6, 64)
	return plots.PlotMPL(data, params, seriesParams, grapherPath)
}
