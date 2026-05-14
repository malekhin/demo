package json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func JsonFrom(path string, output interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	err = json.Unmarshal(b, &output)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}
