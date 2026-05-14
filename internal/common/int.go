package common

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

type Int struct {
	*int
}

func NewInt[T constraints.Integer](v T) Int {
	str := fmt.Sprintf("%d", v)
	i, _ := strconv.Atoi(str)

	return Int{int: &i}
}

func (j Int) PtrInt64() *int64 {
	v := j.int
	if v != nil {
		ptr := int64(*v)
		return &ptr
	}

	return nil
}
