package util

func Contains[V comparable](slice []V, value V) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func SliceSplit[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		panic("chunkSize must be greater than zero")
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
