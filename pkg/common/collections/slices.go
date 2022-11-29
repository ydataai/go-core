package collections

// TODO move to go-core

func Filter[T interface{}](list []T, test func(T) bool) []T {
	ret := []T{}
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func Map[S, D interface{}](list []S, mapTo func(S) D) []D {
	ret := []D{}
	for _, item := range list {
		ret = append(ret, mapTo(item))
	}
	return ret
}

func FilterAndMap[S, D interface{}](list []S, test func(S) bool, mapTo func(S) D) []D {
	ret := []D{}
	for _, item := range list {
		if test(item) {
			ret = append(ret, mapTo(item))
		}
	}
	return ret
}
