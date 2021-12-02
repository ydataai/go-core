package helpers

// FirstStringOfArrayWithFallback is a handle to get the first position of the array with fallback
func FirstStringOfArrayWithFallback(list []string, fallback string) string {
	if len(list) > 0 {
		return list[0]
	}
	return fallback
}
