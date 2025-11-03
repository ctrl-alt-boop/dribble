package data

import (
	"encoding/json"
)

type Json string

func (j *Json) PrettifyJson(jsonString string) string {
	var tmp any
	if err := json.Unmarshal([]byte(jsonString), &tmp); err != nil {
		return jsonString
	}
	prettyJson, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return jsonString
	}
	return string(prettyJson)
}

func (j Json) Render() string {
	return j.PrettifyJson(string(j))
}
