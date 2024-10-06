package utils

import (
	"encoding/json"
)

func StructToMap(input interface{}) (map[string]interface{}, error) {
	// Chuyển struct thành JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Giải mã JSON thành map[string]interface{}
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func GetBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}