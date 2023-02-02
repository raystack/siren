package pgc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringInterfaceMap map[string]interface{}

type StringStringMap map[string]string

func (m *StringInterfaceMap) Scan(value interface{}) error {
	if value == nil {
		m = new(StringInterfaceMap)
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringInterfaceMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

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
