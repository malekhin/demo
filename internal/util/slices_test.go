package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitInt64(t *testing.T) {
	cases := []struct {
		input  []int64
		size   int
		output [][]int64
	}{
		{
			input:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			size:   3,
			output: [][]int64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
		},
		{
			input:  []int64{1, 2, 3},
			size:   3,
			output: [][]int64{{1, 2, 3}},
		},
	}

	for _, c := range cases {
		res := SliceSplit(c.input, c.size)
		assert.Equal(t, c.output, res)
	}
}
