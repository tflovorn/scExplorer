package tempAll

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)
import "../serialize"

func LoadEnvCache(cachePath string) ([]interface{}, []error, error) {
	jsonData, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, nil, err
	}
	cache := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &cache)
	if err != nil {
		return nil, nil, err
	}
	ifData := cache["data"].([]interface{})
	data := make([]interface{}, len(ifData))
	for i, d := range ifData {
		env := new(Environment)
		md := d.(map[string]interface{})
		serialize.CopyValues(&md, env)
		data[i] = *env
	}
	ifErrs := cache["errs"].([]interface{})
	errs := make([]error, len(ifErrs))
	for i, err := range ifErrs {
		strErr := err.(string)
		if strErr != "" {
			errs[i] = errors.New(err.(string))
		}
	}
	return data, errs, nil
}

func SaveEnvCache(cachePath string, data []interface{}, errs []error) error {
	cache := make(map[string]interface{})
	cache["data"] = data
	errStrings := make([]string, len(errs))
	for i, err := range errs {
		if err != nil {
			errStrings[i] = err.Error()
		}
	}
	cache["errs"] = errStrings
	jsonData, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	ioutil.WriteFile(cachePath, jsonData, os.ModePerm)
	return nil
}
