package canvas

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
	ImageType       string // 图片类型,
	Text            string // 文本
	CenterType      string // 居中对齐类型
	WordWrap        string // 自动换行
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

type ImageType uint8

const (
	ImageTypePng  = ImageType(1)
	ImageTypeJpeg = ImageType(2)
)
