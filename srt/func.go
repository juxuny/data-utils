package srt

func splitTagRawData(data []byte) [][]byte {
	splitter := map[byte]bool{
		'<': true,
		'>': true,
		'=': true,
		' ': true,
	}
	ret := make([][]byte, 0)
	buf := make([]byte, 0)
	quote := 0
	for i := range data {
		if splitter[data[i]] && quote == 0 {
			if len(buf) > 0 {
				ret = append(ret, buf)
				buf = make([]byte, 0)
			}
			continue
		}
		if data[i] == '"' && quote == 0 {
			quote += 1
		} else if data[i] == '"' && quote > 0 {
			quote -= 1
		}
		buf = append(buf, data[i])
	}
	if len(buf) > 0 {
		ret = append(ret, buf)
	}
	return ret
}

func Trim(data []byte, charset string) []byte {
	set := map[byte]bool{}
	for i := range charset {
		set[charset[i]] = true
	}
	start := 0
	end := len(data) - 1
	for start <= end {
		if !set[data[start]] {
			break
		}
		start += 1
	}
	for start <= end {
		if !set[data[end]] {
			break
		}
		end -= 1
	}
	ret := make([]byte, end-start+1)
	copy(ret, data[start:])
	return ret
}
