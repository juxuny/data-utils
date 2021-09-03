package lib

type _Int struct{}

var Int = _Int{}

func (_Int) CompareSlice(a, b []int) (isEqual bool) {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
