package util

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](elements ...T) T {
	ret := elements[0]

	for _, x := range elements {
		if x < ret {
			ret = x
		}
	}

	return ret
}

func Max[T constraints.Ordered](elements ...T) T {
	ret := elements[0]

	for _, x := range elements {
		if x > ret {
			ret = x
		}
	}

	return ret
}
