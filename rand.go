package data_utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const tb = "qazxswedcvfrtgbnhyujmkiolpQAZXSWEDCVFRTGBNHYUJMKIOLP1234567890"

func RandString(l int) string {
	buf := make([]byte, l)
	for i := 0; i < l; i++ {
		buf[i] = tb[rand.Intn(len(tb))]
	}
	return string(buf)
}
