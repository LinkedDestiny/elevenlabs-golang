package types

import (
	"encoding/json"
	"reflect"
)

// Validator interface for types that can validate themselves
type Validator interface {
	Validate() error
}

// BaseModel provides common functionality for all model types
type BaseModel struct{}

// RemoveNilFields removes nil fields from a struct for JSON serialization
func RemoveNilFields(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get the JSON tag
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName := fieldType.Name
		if jsonTag != "" {
			fieldName = jsonTag
		}

		// Only include non-nil values
		if !field.IsNil() {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

// ToJSON converts a struct to JSON, removing nil fields
func ToJSON(v interface{}) ([]byte, error) {
	cleaned := RemoveNilFields(v)
	return json.Marshal(cleaned)
}
