package utils

//求uint数组交集
func Intersect(slice1, slice2 []uint32) []uint32 {
	if len(slice1) == 0 {
		return slice2
	}
	if len(slice2) == 0 {
		return slice1
	}
	m := make(map[uint32]uint32)
	n := make([]uint32, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times > 0 {
			n = append(n, v)
		}
	}
	return n
}
