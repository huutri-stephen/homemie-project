package utils

import (
	"encoding/json"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/shopspring/decimal"
)

// ConvertStringArrayToJSON converts a string array to a JSON string.
func ConvertStringArrayToJSON(stringArray []string) string {
	jsonBytes, err := json.Marshal(stringArray)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// FloatToDecimal converts a float64 to a types.Decimal.
func FloatToDecimal(f float64) types.Decimal {
    d := decimal.NewFromFloat(f)
    return types.Decimal{Decimal: d}
}