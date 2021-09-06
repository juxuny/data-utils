package dict

import (
	"fmt"
	"github.com/juxuny/env/ks"
)

type Attr map[string]interface{}

var AttrKey = struct {
	FontFace        string
	FontSize        string
	FontColor       string
	Color           string
	BackgroundColor string
	Src             string
	ImageType       string
	Text            string
}{}

func init() {
	ks.InitKeyName(&AttrKey, false)
}

func invalidAttr(k string, value ...interface{}) error {
	s := fmt.Sprintf("invalid attr:%s", k)
	if len(value) > 0 {
		s += fmt.Sprintf("=%v", value[0])
	}
	return fmt.Errorf(s)
}
