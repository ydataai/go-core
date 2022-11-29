package collections

// Filter an slice given a filter function
func Filter[T interface{}](list []T, test func(T) bool) []T {
	ret := []T{}
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

// Map convert a list given a map function
func Map[S, D interface{}](list []S, mapTo func(S) D) []D {
	ret := []D{}
	for _, item := range list {
		ret = append(ret, mapTo(item))
	}
	return ret
}

// FilterAndMap filter a list and converts to other list given a filter and a map functions
func FilterAndMap[S, D interface{}](list []S, test func(S) bool, mapTo func(S) D) []D {
	ret := []D{}
	for _, item := range list {
		if test(item) {
			ret = append(ret, mapTo(item))
		}
	}
	return ret
}
