package utils

import (
	"encoding/json"

	"github.com/ericlagergren/decimal"
)

// ConvertStringArrayToJSON converts a string array to a JSON string.
func ConvertStringArrayToJSON(stringArray []string) string {
	jsonBytes, err := json.Marshal(stringArray)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func FloatToDecimal(f float64) *decimal.Big {
    return new(decimal.Big).SetFloat64(f)
}