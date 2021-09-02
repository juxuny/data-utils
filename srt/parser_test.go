package srt

import "testing"

func TestParse(t *testing.T) {
	input := `1
00:00:37,700 --> 00:00:41,120
<font face="Cronos Pro Light" size="38" color="#FF0000"><i>Legend tells of a <b><font color="#00FF00">legendary</font></b> warrior</i></font>

2
00:00:41,200 --> 00:00:44,370
<font face="Cronos Pro Light" size="38" color="#fff0cf"><i>whose kung fu skills</i>
<i>were the stuff of legend.</i></font>
`
	blocks, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blocks)
}
