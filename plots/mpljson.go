package plots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const FILE_KEY = "graph_filepath"
const SERIES_KEY = "series"
const DATA_KEY = "data"
const XLABEL_KEY = "xlabel"
const YLABEL_KEY = "ylabel"
const YMIN_KEY = "ymin"

// A map compatible with the JSON object type
type jsonObject map[string]interface{}

// Plot the given data using matplotlib. Global graph parameters are given
// by params; series-specific parameters are given by seriesParams.
// params[FILE_KEY] specifies the path+fileName.
func PlotMPL(data []Series, params map[string]string, seriesParams []map[string]string, grapherPath string) error {
	filePath, err := graphDataToFile(data, params, seriesParams)
	if err != nil {
		return err
	}
	cmd := exec.Command("/usr/bin/env", "python", grapherPath, filePath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Create a seriesParams for PlotMPL
func MakeSeriesParams(varLabel, varFmt string, varVals []float64, styles []string) []map[string]string {
	ret := make([]map[string]string, len(varVals))
	for i, val := range varVals {
		label := fmt.Sprintf("$%v="+varFmt+"$", varLabel, val)
		sp := make(map[string]string)
		sp["label"] = label
		if i < len(styles) {
			sp["style"] = styles[i]
		} else {
			sp["style"] = styles[0]
		}
		ret[i] = sp
	}
	return ret
}

// Write the graph data to the file specified in params[FILE_KEY]
func graphDataToFile(data []Series, params map[string]string, seriesParams []map[string]string) (string, error) {
	// place data in JSON object to be marshalled
	target := jsonObject{}
	for k, v := range params {
		target[k] = v
	}
	target[SERIES_KEY] = make([]jsonObject, len(data))
	for i := 0; i < len(data); i++ {
		sp := jsonObject{}
		for k, v := range seriesParams[i] {
			sp[k] = v
		}
		sp[DATA_KEY] = data[i].Pairs()
		target[SERIES_KEY].([]jsonObject)[i] = sp
	}
	// marshal to JSON
	marshalled, err := json.Marshal(target)
	if err != nil {
		return "", err
	}
	// extract filepath
	filePrefix, ok := params[FILE_KEY]
	if !ok {
		return "", fmt.Errorf("no %v in params", FILE_KEY)
	}
	jsonPath := filePrefix + ".json"
	// write file
	ioutil.WriteFile(jsonPath, marshalled, os.ModePerm)
	return jsonPath, nil
}
