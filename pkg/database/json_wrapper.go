package database

import (
	"database/sql/driver"
	"demo/internal/common"
	"encoding/json"
	"errors"
	"fmt"
)

func NewJsonWrapper(i interface{}) JsonWrapper {
	return JsonWrapper{Data: i}
}

type JsonWrapper struct {
	Data interface{}
}

func (j JsonWrapper) Value() (driver.Value, error) {
	b, err := json.Marshal(j.Data)
	if err != nil {
		return nil, common.Wrap(err, "Marshal")
	}

	s := string(b)
	if s == "null" {
		return "{}", nil
	}

	return s, nil
}

func (j *JsonWrapper) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var err error
	switch v := src.(type) {
	case []byte:
		err = json.Unmarshal(v, &j.Data)
	case string:
		err = json.Unmarshal([]byte(v), &j.Data)
	default:
		return errors.New("type assertion failed")
	}
	if err != nil {
		return common.Wrap(err, "Unmarshal")
	}

	return nil
}

func (j *JsonWrapper) GetInt(key string) (int, error) {
	m, ok := j.Data.(map[string]interface{})
	if !ok {
		return 0, nil
	}

	val, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' is not exists", key)
	}

	res, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("key '%s' is not numeric", key)
	}

	return int(res), nil
}

func (j *JsonWrapper) Get(output interface{}) error {
	m, ok := j.Data.(map[string]interface{})
	if !ok {
		return errors.New("JsonWrapper.Data is not a map")
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return errors.New("json.Marshal")
	}

	err = json.Unmarshal(jsonBytes, output)
	if err != nil {
		return common.Wrap(err, "json.Unmarshal")
	}

	return nil
}

func (j *JsonWrapper) GetJson() (string, error) {
	v, err := j.Value()
	if err != nil {
		return "", common.Wrap(err, "Marshal")
	}

	return v.(string), nil
}
