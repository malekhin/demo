package common

func SliceSplit(slice []int, size int) [][]int {
	if size <= 0 {
		return nil
	}

	var res [][]int
	for len(slice) > 0 {
		if len(slice) < size {
			size = len(slice)
		}
		res = append(res, slice[:size])
		slice = slice[size:]
	}

	return res
}
