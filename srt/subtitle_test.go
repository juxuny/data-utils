package srt

import (
	"testing"
)

func TestSubtitleParse(t *testing.T) {
	data := `<font face="Cronos Pro Light" size="38" color="#FF0000"><i>Legend tells of a <b><font color="#00FF00">legendary</font></b> warrior</i></font>`
	//data = `<f><i>L<b><f>dd</f></b>w</i></f>`
	//data = `<f><i>L<b>a</b><b>a</b>w</i></f>`
	//data = `<f>3434</f>`
	//data = `I'm issuing<font face="Cronos Pro Light" size="38" color="#FF0000"><i>Legend tells of a <b><font color="#00FF00">legendary</font></b> warrior</i></font>I'm issuing `
	node, err := ParseSubtitle([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range node {
		t.Log(item)
	}
}

func TestSplitByCharset(t *testing.T) {
	data := `<f size="abc"   family="Chas 1344">`
	l := splitTagRawData([]byte(data))
	for _, buf := range l {
		t.Log(string(Trim(buf, " ")))
	}
}
