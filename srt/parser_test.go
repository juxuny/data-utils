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

12
00:01:27,330 --> 00:01:30,090
<font face="Cronos Pro Light" size="38" color="#fff0cf">-And attractive.
-How can we repay you?</font>
`
	blocks, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range blocks {
		t.Log(b.String())
	}
}

func TestParseInterval(t *testing.T) {
	input := "00:00:37,700 --> 00:00:41,120"
	start, end, err := parseInterval(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(start.Format("15:04:05.000"))
	t.Log(end.Format("15:04:05.000"))
}
