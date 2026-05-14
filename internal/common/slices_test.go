package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceSplit(t *testing.T) {
	tests := []struct {
		Slice    []int
		Size     int
		Expected [][]int
	}{
		{
			Slice:    []int{1, 2, 3},
			Size:     1,
			Expected: [][]int{{1}, {2}, {3}},
		},
		{
			Slice:    []int{1, 2, 3},
			Size:     2,
			Expected: [][]int{{1, 2}, {3}},
		},
		{
			Slice:    []int{1, 2, 3},
			Size:     3,
			Expected: [][]int{{1, 2, 3}},
		},
		{
			Slice:    []int{1, 2, 3},
			Size:     4,
			Expected: [][]int{{1, 2, 3}},
		},
		{
			Slice:    []int{1, 2, 3},
			Size:     0,
			Expected: nil,
		},
	}

	for _, test := range tests {
		res := SliceSplit(test.Slice, test.Size)
		assert.Equal(t, test.Expected, res)
	}
}
