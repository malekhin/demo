package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Int(t *testing.T) {
	test := 1
	i := NewInt(test)
	res := i.PtrInt64()
	assert.Equal(t, int(*res), test)
}
