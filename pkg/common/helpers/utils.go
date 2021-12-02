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

// FirstStringOfArrayWithFallback is a handle to get the first position of the array with fallback
func FirstStringOfArrayWithFallback(list []string, fallback ...string) string {
	if len(list) > 0 {
		return list[0]
	} else if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// LastStringOfArrayWithFallback is a handle to get the last position of the array with fallback
func LastStringOfArrayWithFallback(list []string, fallback ...string) string {
	if len(list) > 0 {
		return list[len(list)-1]
	} else if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
