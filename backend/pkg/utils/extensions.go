package utils

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StructToMap(input interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

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

func GetArray(data map[string]interface{}, key string) []interface{} {
	if val, ok := data[key]; ok {
		if array, ok := val.([]interface{}); ok {
			return array
		}
	}
	return []interface{}{}
}

func GetTime(data map[string]interface{}, key string) time.Time {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case time.Time:
			return v
		case string:
			parsedTime, err := time.Parse(time.RFC3339, v)
			if err == nil {
				return parsedTime
			}
		}
	}
	return time.Time{}
}

func GetStringArray(data map[string]interface{}, key string) []string {
	if val, ok := data[key]; ok {
		if array, ok := val.([]interface{}); ok {
			strArray := make([]string, len(array))
			for i, v := range array {
				if str, ok := v.(string); ok {
					strArray[i] = str
				}
			}
			return strArray
		}
	}
	return []string{}
}

func GetMap(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal
		}
	}
	return map[string]interface{}{}
}

func GetObjectID(data map[string]interface{}, key string) primitive.ObjectID {
	if id, ok := data[key]; ok {
		switch v := id.(type) {
		case primitive.ObjectID:
			return v
		case string:
			objectID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return primitive.NilObjectID
			}
			return objectID
		}
	}
	return primitive.NilObjectID
}
