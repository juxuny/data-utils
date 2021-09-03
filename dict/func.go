package dict

import (
	"github.com/juxuny/data-utils/lib"
	"unicode"
)

func ReplaceWords(input []byte, oldStr []byte, newStr []byte) []byte {
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
	length := 0
	for i := 0; i < len(list); i++ {
		if lib.Byte.CompareSlice(list[i], oldStr) {
			list[i] = newStr
			length += len(newStr)
		} else {
			length += len(list[i])
		}
	}
	ret := make([]byte, length)
	i = 0
	for _, l := range list {
		copy(ret[i:(i+len(l))], l)
		i += len(l)
	}

	return ret
}
