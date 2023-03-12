package utils

func Uniques(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func Contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}
