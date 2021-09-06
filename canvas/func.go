package dict

import (
	"encoding/hex"
	"image/color"
	"strings"
)

func convertHexToRGBA(c string) (color.RGBA, error) {
	ret := color.RGBA{A: 0xFF}
	s := strings.Trim(c, "# ")
	b, err := hex.DecodeString(s)
	if err != nil {
		return ret, err
	}
	if len(s) == 8 {
		ret.A = b[0]
		ret.R = b[1]
		ret.G = b[2]
		ret.B = b[3]
	} else {
		ret.R = b[0]
		ret.G = b[1]
		ret.B = b[2]
	}
	return ret, nil
}
