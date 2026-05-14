package converter

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func StructToJson(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("%w: Marshal", err)
	}

	return string(b), nil
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	err := mapstructure.Decode(input, &output)
	if err != nil {
		return fmt.Errorf("%w: Decode", err)
	}

	return nil
}

func StructToMap(input interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%w: Marshal", err)
	}
	var output map[string]interface{}
	err = json.Unmarshal(b, &output)
	if err != nil {
		return nil, fmt.Errorf("%w: Unmarshal", err)
	}

	return output, nil
}
