package helper

import (
	"errors"
	"fmt"
)

func GetMapString(m map[string]interface{}, name string, key string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("No value supplied for required %s map key %q", name, key))
	}
	typedVal, ok := val.(string)
	if !ok {
		return "", errors.New(
			fmt.Sprintf(
				"Wrong type for %s map key %q: expected type %v, got value %q of type %t",
				name, key, "string", val, val,
			),
		)
	}
	return typedVal, nil
}
