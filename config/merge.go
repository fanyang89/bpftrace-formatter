package config

import (
	"encoding/json"
	"fmt"
)

// MergeConfig deep-merges overrides into base. Unknown keys are ignored with warnings.
func MergeConfig(base *Config, overrides ...map[string]any) (*Config, []string, error) {
	if base == nil {
		base = DefaultConfig()
	}

	baseMap, err := configToMap(base)
	if err != nil {
		return nil, nil, err
	}

	warnings := []string{}
	for _, override := range overrides {
		if override == nil {
			continue
		}
		baseMap = mergeMap(baseMap, override, "", &warnings)
	}

	merged, err := mapToConfig(baseMap)
	if err != nil {
		return nil, warnings, err
	}

	return merged, warnings, nil
}

func configToMap(cfg *Config) (map[string]any, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func mapToConfig(value map[string]any) (*Config, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func mergeMap(base map[string]any, override map[string]any, path string, warnings *[]string) map[string]any {
	for key, overrideValue := range override {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		baseValue, exists := base[key]
		if !exists {
			*warnings = append(*warnings, fmt.Sprintf("unknown config key: %s", currentPath))
			continue
		}

		baseMap, baseIsMap := baseValue.(map[string]any)
		overrideMap, overrideIsMap := overrideValue.(map[string]any)
		if baseIsMap && overrideIsMap {
			base[key] = mergeMap(baseMap, overrideMap, currentPath, warnings)
			continue
		}

		base[key] = overrideValue
	}

	return base
}
