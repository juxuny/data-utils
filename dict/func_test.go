package dict

import (
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
