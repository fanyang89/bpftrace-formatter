package config

import (
	"fmt"
	"reflect"
	"strings"
)

// MergeConfig deep-merges overrides into base. Unknown keys are ignored with warnings.
// It uses reflection to provide a type-safe and efficient merge without JSON marshaling.
func MergeConfig(base *Config, overrides ...map[string]any) (*Config, []string, error) {
	if base == nil {
		base = DefaultConfig()
	}

	// Work on a copy of the base config
	merged := *base
	warnings := []string{}

	for _, override := range overrides {
		if override == nil {
			continue
		}
		mergeStructWithMap(reflect.ValueOf(&merged).Elem(), override, "", &warnings)
	}

	return &merged, warnings, nil
}

// mergeStructWithMap recursively merges a map into a struct using reflection.
func mergeStructWithMap(structVal reflect.Value, override map[string]any, path string, warnings *[]string) {
	structType := structVal.Type()

	for key, overrideValue := range override {
		fieldName, found := findFieldByJSONTag(structType, key)
		if !found {
			fullPath := key
			if path != "" {
				fullPath = path + "." + key
			}
			*warnings = append(*warnings, fmt.Sprintf("unknown config key: %s", fullPath))
			continue
		}

		fieldVal := structVal.FieldByName(fieldName)
		if !fieldVal.IsValid() || !fieldVal.CanSet() {
			continue
		}

		fullPath := key
		if path != "" {
			fullPath = path + "." + key
		}

		// Handle nested structs
		if fieldVal.Kind() == reflect.Struct {
			if overrideMap, ok := overrideValue.(map[string]any); ok {
				mergeStructWithMap(fieldVal, overrideMap, fullPath, warnings)
				continue
			}
		}

		// Set the value directly
		setFieldValue(fieldVal, overrideValue, fullPath, warnings)
	}
}

// findFieldByJSONTag finds a struct field name by its JSON tag.
func findFieldByJSONTag(t reflect.Type, tagValue string) (string, bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		// JSON tags can have options like "name,omitempty"
		parts := strings.Split(tag, ",")
		if parts[0] == tagValue {
			return field.Name, true
		}
	}
	return "", false
}

// setFieldValue attempts to set a reflect.Value from an interface{} with type conversion.
func setFieldValue(fieldVal reflect.Value, val any, path string, warnings *[]string) {
	if val == nil {
		return
	}

	valVal := reflect.ValueOf(val)

	// Direct assignment if types match
	if valVal.Type().AssignableTo(fieldVal.Type()) {
		fieldVal.Set(valVal)
		return
	}

	// Handle numeric conversions (JSON numbers are often float64)
	if valVal.Kind() == reflect.Float64 {
		switch fieldVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldVal.SetInt(int64(valVal.Float()))
			return
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldVal.SetUint(uint64(valVal.Float()))
			return
		case reflect.Float32, reflect.Float64:
			fieldVal.SetFloat(valVal.Float())
			return
		}
	}

	// If we're here, types don't match and couldn't be easily converted
	*warnings = append(*warnings, fmt.Sprintf("invalid value type for %s: expected %v, got %v", path, fieldVal.Type(), valVal.Type()))
}
