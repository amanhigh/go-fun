package util

func SliceContains(key string, slice []string) bool {
	for _, value := range slice {
		if key == value {
			return true
		}
	}
	return false
}
