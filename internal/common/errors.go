package common

import "fmt"

func Wrap(err error, message string) error {
	return fmt.Errorf("%v: %w", message, err)
}
