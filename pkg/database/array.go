package database

import (
	"database/sql/driver"
	"demo/internal/common"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type IntArr []int
type StringArr []string

func (i *IntArr) Scan(src interface{}) error {
	var err error
	switch v := src.(type) {
	case []byte:
		s := strings.Replace(string(v), `"`, "", -1)
		err = json.Unmarshal([]byte(s), i)
	default:
		return errors.New("type assertion failed")
	}
	if err != nil {
		return common.Wrap(err, "Unmarshal")
	}

	return nil
}

func (i IntArr) Value() (driver.Value, error) {
	j := make([]string, 0, len(i))
	for _, n := range i {
		j = append(j, fmt.Sprintf("%d", n))
	}

	return fmt.Sprintf("{%s}", strings.Join(j, ",")), nil
}

func (s *StringArr) Scan(src interface{}) error {
	var err error
	switch v := src.(type) {
	case []byte:
		err = json.Unmarshal(v, s)
	case string:
		err = json.Unmarshal([]byte(v), s)
	default:
		return errors.New("type assertion failed")
	}
	if err != nil {
		return common.Wrap(err, "Unmarshal")
	}

	return nil
}

func (j StringArr) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(j, ",")), nil
}
