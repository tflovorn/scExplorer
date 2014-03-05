import sys
import json
import math
import matplotlib.pyplot as plt
from matplotlib.ticker import FormatStrFormatter
from matplotlib.font_manager import FontProperties
from numpy import arange, meshgrid

_GRAPH_DEFAULTS = {"xlabel":"$x$", "ylabel":"$y$", "num_ticks":5, 
    "axis_label_fontsize":"x-large", "tick_formatstr":"%.2f",
    "legend_fontsize":"large", "legend_loc":0, "legend_title":None, 
    "ymin":None, "graph_filepath":None, "plot_type": "scatter",
    "th":None,"thp":None,"t0":None,"D1":None,"Mu_h":None,"epsilon_min":None}

_SERIES_DEFAULTS = {"label":None, "style":"k."}

def parse_file(file_path):
    '''Return the plot representation of the JSON file specified.'''
    # -- todo : check for IOError --
    return import_json(open(file_path, 'r').read())

def import_json(json_string):
    '''Return the plot representation of the given JSON string.'''
    graph_data = json.loads(json_string)
    if isinstance(graph_data, list):
        graph_data = [add_default_data(graph) for graph in graph_data]
    else:
        graph_data = add_default_data(graph_data)
    return graph_data

def add_default_data(graph_data):
    # graph-wide defaults
    for key, value in _GRAPH_DEFAULTS.items():
        if key not in graph_data:
            graph_data[key] = value
    # hack to give a fresh series list each time
    if "series" not in graph_data:
        graph_data["series"] = []
    # series-specific defaults
    for series in graph_data["series"]:
        for key, value in _SERIES_DEFAULTS.items():
            if key not in series:
                series[key] = value
    return graph_data

def make_graph(graph_data):
    '''Take a dictionary representing a graph or a list of such dictionaries.
    Build the graph(s), save them to file(s) (if requested), and return the
    matplotlib figures.

    '''
    # Process a list of graphs one element at a time.
    if isinstance(graph_data, list):
        return [make_graph(some_graph) for some_graph in graph_data]
    # If we need to make a Fermi surface plot, go to the function for that.
    if graph_data["plot_type"] == "Fermi_surface":
        return plot_Fermi_surface(graph_data)
    # If we get here, this is not a Fermi surface plot - make scatter plot instead.
    try:
        dims = graph_data["dimensions"]
        fig = plt.figure(figsize=(dims[0], dims[1]))
    except:
        fig = plt.figure()
    axes = fig.add_subplot(1, 1, 1)
    bounds = [None, None]
	# plot the data
    for series in graph_data["series"]:
        fig,axes,bounds = _graph_series(graph_data,series,fig,axes,bounds)
	# set properties
    fontprop_legend = FontProperties(size=graph_data["legend_fontsize"])
    axes.legend(loc=graph_data["legend_loc"], title=graph_data["legend_title"],
                prop=fontprop_legend)
    axes.set_xlabel(graph_data["xlabel"], size=graph_data["axis_label_fontsize"])
    axes.set_ylabel(graph_data["ylabel"], size=graph_data["axis_label_fontsize"])
    if graph_data["ymin"] != None and graph_data["ymin"] != "":
        axes.set_ylim(bottom=float(graph_data["ymin"]), auto=None)
    _save_figure(graph_data, fig)
    return fig, axes

def _graph_series(graph_data, series, fig, axes, bounds):
    # -- todo : set ticks --
    axes.plot(_xData(series), _yData(series), series["style"], 
              label=series["label"])
    return fig, axes, bounds

def _xData(series):
    return [point[0] for point in series["data"]]

def _yData(series):
    return [point[1] for point in series["data"]]

def _save_figure(graph_data, fig):
    if graph_data["graph_filepath"] is None:
        return
    fig.savefig(graph_data["graph_filepath"] + ".png")
    fig.savefig(graph_data["graph_filepath"] + ".eps")

# Plot a single Fermi surface.
# To make this plot, graph_data["Fermi_surface_data"] must be a dictionary
# with the keys "th", "thp", "t0", "D1", "Mu_h", and "epsilon_min".
def plot_Fermi_surface(graph_data):
    delta = 0.05
    x = arange(-math.pi, math.pi, delta)
    y = arange(-math.pi, math.pi, delta)
    X, Y = meshgrid(x, y)
    FS = lambda x, y: _step(_xi_h(graph_data, x, y))
    Z = []
    for i in range(len(y)):
        Z.append([])
        for j in range(len(x)):
            Z[i].append(FS(x[j], y[i]))
    fig = plt.figure()
    axes = fig.add_subplot(1, 1, 1)
    CS = axes.contour(X, Y, Z)
    _save_figure(graph_data, fig)

# if epsilon_min is included, step(xi) = 1 for all k-points - why?
def _xi_h(fsd, kx, ky):
    sx, sy = math.sin(kx), math.sin(ky)
    envVars = map(float, [fsd["th"], fsd["D1"], fsd["t0"], fsd["thp"], fsd["epsilon_min"], fsd["Mu_h"]])
    th, D1, t0, thp, epsilon_min, Mu_h = envVars
    #eps = 2.0*th*((sx+sy)*(sx+sy) - 1.0) + 4.0*(2.0*D1*t0 - thp)*sx*sy - epsilon_min
    eps = 2.0*th*((sx+sy)*(sx+sy) - 1.0) + 4.0*(2.0*D1*t0 - thp)*sx*sy
    return eps - Mu_h

def _step(x):
    if x < 0.0:
        return 0.0
    else:
        return 1.0

if __name__ == "__main__":
    if len(sys.argv) > 1:
        make_graph(parse_file(sys.argv[1]))
