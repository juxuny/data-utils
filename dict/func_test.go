package dict

import (
	"fmt"
	"testing"
)

func TestReplaceAllWords(t *testing.T) {
	var dataList = []struct {
		Input  string
		Result string
		Old    string
		New    string
	}{
		{Input: "to be or not to be", Old: "be", New: "me", Result: "to me or not to me"},
	}
	for _, d := range dataList {
		result := ReplaceWords([]byte(d.Input), []byte(d.Old), []byte(d.New))
		if string(result) != d.Result {
			t.Fatal("wrong result: ", string(result))
		}
	}
}

func TestConvertHexColor(t *testing.T) {
	s := "#AA0F0E"
	c, err := convertHexToRGBA(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%02x %02x %02x %02x", c.A, c.R, c.G, c.B))
}
