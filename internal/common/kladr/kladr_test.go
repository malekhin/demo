package kladr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertKladrIdsMapsRegionIdsToKladrCodes(t *testing.T) {
	tests := []struct {
		regionIds []int
		kladrIds  []string
	}{
		{
			regionIds: []int{83},
			kladrIds:  []string{"20"},
		},
		{
			regionIds: []int{18, 78},
			kladrIds:  []string{"18", "79"},
		},
		{
			regionIds: []int{79, 80, 81, 82, 84, 85, 87},
			kladrIds:  []string{"83", "86", "87", "89", "91", "92", "99"},
		},
		{
			regionIds: []int{100},
			kladrIds:  []string{"100"},
		},
	}

	for _, tt := range tests {
		result := ConvertKladrIds(tt.regionIds)

		assert.Len(t, result, len(tt.kladrIds))
		assert.Equal(t, tt.kladrIds, result)
	}
}
