package kladr

import (
	"fmt"
	"strconv"
)

func ConvertKladrIds(regionIds []int) []string {
	kladrIds := make([]string, len(regionIds))
	for i, regionId := range regionIds {
		kladrIds[i] = ConvertKladrId(regionId)
	}

	return kladrIds
}

func ConvertKladrId(regionId int) string {
	switch true {
	case regionId == 83:
		return "20"
	case regionId < 20:
		return fmt.Sprintf("%02s", strconv.Itoa(regionId))
	case regionId < 79:
		return strconv.Itoa(regionId + 1)
	case regionId == 79:
		return "83"
	case regionId == 80:
		return "86"
	case regionId == 81:
		return "87"
	case regionId == 82:
		return "89"
	case regionId == 84:
		return "91"
	case regionId == 85:
		return "92"
	case regionId == 87:
		return "99"
	default:
		return strconv.Itoa(regionId)
	}
}
