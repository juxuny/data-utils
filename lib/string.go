package lib

import "unicode"

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

type _String struct{}

var String = _String{}

func (_String) SplitWithStopFunc(text string, stopFunc func(x rune) bool) []string {
	input := []byte(text)
	if len(input) == 0 {
		return []string{}
	}
	var list [][]byte
	buf := []byte{input[0]}
	i := 1
	for i < len(input) {
		if unicode.IsLetter(rune(input[i])) != unicode.IsLetter(rune(input[i-1])) {
			list = append(list, buf)
			buf = []byte{input[i]}
		} else {
			buf = append(buf, input[i])
		}
		i++
	}
	if len(buf) > 0 {
		list = append(list, buf)
	}
	ret := make([]string, 0)
	for _, item := range list {
		ret = append(ret, string(item))
	}
	return ret
}
