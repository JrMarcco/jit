package slice

func shrink(cap int, length int) (int, bool) {
	if length == 0 || cap == length {
		return cap, false
	}

	// 计算容量与长度的比例
	radio := float32(cap) / float32(length)

	switch {
	// 超大容量：当比例 >= 2 时缩容到 1.5 倍
	case cap > 4096 && radio >= 2:
		return int(float32(length) * 1.5), true
	// 大容量：当比例 >= 2 时缩容到原容量的 50%
	case cap > 1024 && radio >= 2:
		return cap / 2, true
	// 中容量：当比例 >= 2.5 时缩容到原容量的 62.5%
	case cap > 256 && radio >= 2.5:
		return int(float32(cap) * 0.625), true
	// 小容量：当比例 >= 3 时缩容到原容量的 50%
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
