package common

func InArray[T int | int64 | string](val T, list []T) bool {
	for _, item := range list {
		if val == item {
			return true
		}
	}

	return false
}
