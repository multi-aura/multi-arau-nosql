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
