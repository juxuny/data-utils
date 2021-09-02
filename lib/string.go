package lib

func StringSliceFilter(in []string, f func(item string) bool) []string {
	ret := make([]string, 0)
	for _, item := range in {
		if f(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

type StringSlice []string

func (t StringSlice) Filter(f func(item string) bool) StringSlice {
	ret := make(StringSlice, 0)
	for _, item := range t {
		if f(item) {
			ret = append(ret, item)
		}
	}
	return ret
}
