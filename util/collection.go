package util

func SliceContains(key string, slice []string) bool {
	for _, value := range slice {
		if key == value {
			return true
		}
	}
	return false
}

func SliceMinus(mainSet []string, minusSet []string) (resultSet []string) {
	removeMap := map[string]bool{}
	for _, value := range minusSet {
		removeMap[value] = true
	}

	for _, value := range mainSet {
		if _, ok := removeMap[value]; !ok {
			resultSet = append(resultSet, value)
		}
	}
	return
}
