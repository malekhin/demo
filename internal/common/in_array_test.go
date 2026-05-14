package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testIntCases = []struct {
	params  []int
	val     int
	inArray bool
}{
	{
		params:  []int{1, 2, 3},
		val:     0,
		inArray: false,
	},
	{
		params:  []int{1, 2, 3},
		val:     1,
		inArray: true,
	},
	{
		params:  []int{},
		val:     0,
		inArray: false,
	},
}

var testStringCases = []struct {
	params  []string
	val     string
	inArray bool
}{
	{
		params:  []string{"a", "b", "c"},
		val:     "a",
		inArray: true,
	},
	{
		params:  []string{"a", "b", "c"},
		val:     "d",
		inArray: false,
	},
}

func Test_InArray(t *testing.T) {
	for _, testCase := range testIntCases {
		assert.Equal(t, InArray(testCase.val, testCase.params), testCase.inArray)
	}

	for _, testCase := range testStringCases {
		assert.Equal(t, InArray(testCase.val, testCase.params), testCase.inArray)
	}
}
