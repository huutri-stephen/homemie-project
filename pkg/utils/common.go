package utils

import (
	"encoding/json"
)

// ConvertStringArrayToJSON converts a string array to a JSON string.
func ConvertStringArrayToJSON(stringArray []string) string {
	jsonBytes, err := json.Marshal(stringArray)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
