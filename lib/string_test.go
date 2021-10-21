package lib

import (
	"testing"
	"unicode"
)

func Test_String_SplitWithStopFunc(t *testing.T) {
	data := "Hello World !!! My name is Juxuny Wu."
	l := String.SplitWithStopFunc(data, unicode.IsLetter)
	t.Log(l)
}

func Test_String_TrimSubString(t *testing.T) {
	var data = []struct {
		In  string
		Out string
		Sub string
	}{
		{In: "xxx.subtitle", Out: "xxx", Sub: ".subtitle"},
		{In: "ABCCBA", Out: "BCCB", Sub: "A"},
	}
	for _, d := range data {
		out := String.TrimSubString(d.In, d.Sub)
		if out != d.Out {
			t.Fatal("wrong out: ", out)
		} else {
			t.Log(out)
		}
	}
}
