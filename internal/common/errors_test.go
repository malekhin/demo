package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrDuplicate = errors.New("duplicate")

func Test_Warp(t *testing.T) {
	errTest := Wrap(ErrDuplicate, "message")
	resErr := errors.New("message: duplicate")

	assert.Equal(t, errTest.Error(), resErr.Error())
	assert.ErrorIs(t, errTest, ErrDuplicate)
}
