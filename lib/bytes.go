package lib

type _Byte struct{}

var Byte = _Byte{}

func (_Byte) CompareSlice(a, b []byte) (isEqual bool) {
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

func Contains(data []byte, charset string) bool {
	set := map[byte]bool{}
	for i := range charset {
		set[charset[i]] = true
	}
	for i := range data {
		if set[data[i]] {
			return true
		}
	}
	return false
}

// charset 里任意一个字符都会分割
func SplitByCharset(data []byte, charset string) [][]byte {
	splitter := map[byte]bool{}
	for i := range charset {
		splitter[charset[i]] = true
	}
	ret := make([][]byte, 0)
	buf := make([]byte, 0)
	for i := 0; i < len(data); i++ {
		if splitter[data[i]] {
			if len(buf) > 0 {
				ret = append(ret, buf)
				buf = make([]byte, 0)
			}
			continue
		}
		buf = append(buf, data[i])
	}
	if len(buf) > 0 {
		ret = append(ret, buf)
	}
	return ret
}

func IsQuoted(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return data[0] == '"' && data[len(data)-1] == '"'
}

func (_Byte) IsQuoted(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return data[0] == '"' && data[len(data)-1] == '"'
}

func (_Byte) Trim(data []byte, cut []byte) (ret []byte) {
	s := 0
	for {
		if s+len(cut) > len(data) {
			break
		}
		if Byte.CompareSlice(data[s:(s+len(cut))], cut) {
			s += len(cut)
		} else {
			break
		}
	}
	end := len(data)
	for {
		if end <= 0 {
			end = 0
			break
		}
		if Byte.CompareSlice(data[end-len(cut):end], cut) {
			end -= len(cut)
		} else {
			break
		}
	}
	if s >= end {
		return []byte{}
	}
	return data[s:end]
}

func (_Byte) Drop(data []byte, filter func(r rune) bool) []byte {
	ret := make([]byte, 0)
	for _, x := range data {
		if filter(rune(x)) {
			ret = append(ret, x)
		}
	}
	return ret
}
