package util

import (
	"encoding/json"

	"github.com/muesli/reflow/ansi"
)

func TruncateWithSuffix(str string, num int, suffix string) string {
	if len(str) > num {
		return str[:num-ansi.PrintableRuneWidth(suffix)] + suffix
	}
	return str
}

func Truncate(str string, num int) string {
	if len(str) > num {
		return str[:num]
	}
	return str
}

func PrettifyJson(jsonString string) string {
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
