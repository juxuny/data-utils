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
