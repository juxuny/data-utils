package lib

import (
	"fmt"
	"testing"
)

func Test_Byte_Trim(t *testing.T) {
	var input = []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x1, 0x2}
	var realResult = []byte{0x3, 0x4, 0x5}
	result := Byte.Trim(input, []byte{0x1, 0x2})
	if !Byte.CompareSlice(result, realResult) {
		t.Fatal("failed")
	}
	t.Log(fmt.Sprintf("%02x", result))
}
