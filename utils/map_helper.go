package utils

import (
	"fmt"
)

func GetMapString(m map[string]interface{}, name string, key string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", fmt.Errorf("no value supplied for required %s map key %q", name, key)
	}
	typedVal, ok := val.(string)
	if !ok {
		return "",
			fmt.Errorf(
				"wrong type for %s map key %q: expected type %v, got value %q of type %t",
				name, key, "string", val, val)
	}
	return typedVal, nil
}
