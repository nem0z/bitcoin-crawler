package utils

import (
	"encoding/json"
	"io/ioutil"
)

func ExportMap(path string, data map[string]bool) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ImportMap(path string) (map[string]bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]bool)

	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
