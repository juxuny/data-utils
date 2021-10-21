package lib

import (
	"strings"
	"unicode"
)

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

func (_String) TrimSubString(in, subString string) (out string) {
	return String.TrimSubStringLeft(String.TrimSubStringRight(in, subString), subString)
}

func (_String) Reverse(in string) string {
	b := []byte(in)
	out := make([]byte, len(b))
	copy(out, b)
	for i := 0; i < len(b)>>1; i++ {
		out[i], out[len(out)-1-i] = out[len(out)-1-i], out[i]
	}
	return string(out)
}

func (_String) TrimSubStringLeft(in, subString string) string {
	if len(in) < len(subString) {
		return in
	}
	var out string
	last := in
	out = strings.Replace(in, subString, "", 1)
	for out != last {
		last = out
		out = strings.Replace(in, subString, "", 1)
	}
	return out
}

func (_String) TrimSubStringRight(in, subString string) string {
	if len(in) < len(subString) {
		return in
	}
	reversedInput := String.Reverse(in)
	reversedSubString := String.Reverse(subString)
	return String.Reverse(String.TrimSubStringLeft(reversedInput, reversedSubString))
}
