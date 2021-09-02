package lib

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
