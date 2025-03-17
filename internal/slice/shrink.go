package slice

func shrink(cap int, length int) (int, bool) {
	if length == 0 || cap == length {
		return cap, false
	}

	// calculate the ratio of capacity to length
	radio := float32(cap) / float32(length)

	switch {
	// huge capacity: when the ratio >= 2, shrink to 1.5 times of the original capacity
	case cap > 4096 && radio >= 2:
		return int(float32(length) * 1.5), true
	// large capacity: when the ratio >= 2, shrink to 50% of the original capacity
	case cap > 1024 && radio >= 2:
		return cap / 2, true
	// medium capacity: when the ratio >= 2.5, shrink to 62.5% of the original capacity
	case cap > 256 && radio >= 2.5:
		return int(float32(cap) * 0.625), true
	// small capacity: when the ratio >= 3, shrink to 50% of the original capacity
	case radio >= 3:
		return cap / 2, true
	}

	return cap, false
}

func Shrink[T any](slice []T) []T {
	cap, length := cap(slice), len(slice)

	newCap, shrunken := shrink(cap, length)
	if !shrunken {
		return slice
	}

	res := make([]T, 0, newCap)
	res = append(res, slice...)

	return res
}
