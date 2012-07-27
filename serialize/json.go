package serialize

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Copy values from the JSON string given into object.
func CopyFromJSON(jsonData string, object interface{}) error {
	jsonObject, err := readJSONString(jsonData)
	if err != nil {
		return err
	}
	copyValues(jsonObject, object)
	return nil
}

// Get a JSON object from the string given.
func readJSONString(jsonData string) (*map[string]interface{}, error) {
	jsonBytes, err := stringToBytes(jsonData)
	if err != nil {
		return nil, err
	}
	return readJSONBytes(jsonBytes)
}

// Get a JSON object from the byte slice given.
func readJSONBytes(jsonData []byte) (*map[string]interface{}, error) {
	jsonObject := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &jsonObject)
	if err != nil {
		return nil, err
	}
	return &jsonObject, nil
}

// Serialize object to JSON representation.
func MakeJSON(object interface{}) (string, error) {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}

// Look at each key in jsonObject and copy thats key's value into the
// corresponding field in object.
func copyValues(jsonObject *map[string]interface{}, object interface{}) {
	// dereference the object pointer
	objectValue := reflect.Indirect(reflect.ValueOf(object))
	// iterate over all fields in the JSON object
	for key, value := range *jsonObject {
		// get a reference to the field in object
		field := objectValue.FieldByName(key)
		if !field.CanSet() {
			// Can't set, probably because this field doesn't
			// exist in object.  Skip it silently.
			continue
		}
		// recognize some numeric types which aren't available in JSON
		// (can extend this list)
		fieldType := field.Type().Name()
		if fieldType == "int" {
			value = int(value.(float64))
		} else if fieldType == "uint" {
			value = uint(value.(float64))
		}
		// set the field in object
		field.Set(reflect.ValueOf(value))
	}
}

// Convert string to byte slice
func stringToBytes(str string) ([]byte, error) {
	reader := strings.NewReader(str)
	bytes := make([]byte, len(str))
	for seen := 0; seen < len(str); {
		n, err := reader.Read(bytes)
		if err != nil {
			return nil, err
		}
		seen += n
	}
	return bytes, nil
}
