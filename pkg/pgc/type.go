package pgc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type StringAnyMap map[string]any

func (m *StringAnyMap) Scan(value any) error {
	if value == nil {
		m = new(StringAnyMap)
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringAnyMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

type StringStringMap map[string]string

func (m *StringStringMap) Scan(value interface{}) error {
	if value == nil {
		m = new(StringStringMap)
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringStringMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

type TimeDuration time.Duration

func (t *TimeDuration) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	valueStr, ok := value.(string)
	if !ok {
		return errors.New("failed type assertion to string")
	}

	dur, err := time.ParseDuration(valueStr)
	if err != nil {
		return fmt.Errorf("error parsing duration '%s': %w", valueStr, err)
	}

	*t = TimeDuration(dur)

	return nil
}

func (t TimeDuration) Value() (driver.Value, error) {
	return time.Duration(t).String(), nil
}

type ListStringStringMap []map[string]string

func (m *ListStringStringMap) Scan(value interface{}) error {
	if value == nil {
		m = new(ListStringStringMap)
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a ListStringStringMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}
