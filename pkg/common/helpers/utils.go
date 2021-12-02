package helpers

// ArrayContainsString is a helper to find out if a given string is inside an array or not
func ArrayContainsString(list []string, key string) bool {
	for _, value := range list {
		if key == value {
			return true
		}
	}
	return false
}

// FirstStringOfArray is a handle to get the first position of the array
func FirstStringOfArray(list []string) string {
	if len(list) > 0 {
		return list[0]
	}
	return ""
}

// LastStringOfArray is a handle to get the last position of the array
func LastStringOfArray(list []string) string {
	if len(list) > 0 {
		return list[len(list)-1]
	}
	return ""
}
