package dict

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/pkg/errors"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"strings"
)

type TextView struct {
	*BaseView
	measureError error
	size         int
	src          string
	fontColor    string
	painter      *Painter

	background *image.RGBA
}

func (t *TextView) Draw(img *image.RGBA, vector ...image.Point) error {
	if t.measureError != nil {
		return t.measureError
	}
	start := t.Rect.Min
	if len(vector) > 0 {
		start.X += vector[0].X
		start.Y += vector[0].Y
	}
	draw.Draw(img, image.Rect(start.X, start.Y, start.X+t.Rect.Dx(), start.Y+t.Rect.Dy()), t.background, image.ZP, draw.Src)
	return nil
}

func (t *TextView) Measure() image.Rectangle {
	if t.measured {
		return t.Rect
	}
	text, isString := t.Attr[AttrKey.Text].(string)
	if !isString {
		t.measureError = invalidAttr(AttrKey.Text, t.Attr[AttrKey.Text])
		return t.Rect
	}
	src, isString := t.Attr[AttrKey.Src].(string)
	if !isString {
		t.measureError = invalidAttr(AttrKey.Src, t.Attr[AttrKey.Src])
		return t.Rect
	}
	t.src = src
	size, isInt := t.Attr[AttrKey.FontSize].(int)
	if !isInt {
		t.measureError = fmt.Errorf("%s is not an int", AttrKey.FontSize)
		return t.Rect
	}
	t.size = size
	fontColor, isString := t.Attr[AttrKey.FontColor].(string)
	if !isString {
		t.measureError = invalidAttr(AttrKey.FontColor, t.Attr[AttrKey.FontColor])
		return t.Rect
	}
	t.fontColor = fontColor
	fontData, err := loadFont(src)
	if err != nil {
		t.measureError = fmt.Errorf("measure text-view failed, invalid font data")
		return t.Rect
	}
	t.painter = NewPainter()
	if err := t.painter.SetFontByFontData(fontData); err != nil {
		t.measureError = errors.Wrap(err, "set font failed")
		return t.Rect
	}
	if c, err := convertHexToRGBA(t.fontColor); err != nil {
		t.measureError = invalidAttr(AttrKey.FontColor, t.fontColor)
		return t.Rect
	} else {
		t.painter.SetColor(image.NewUniform(c))
	}
	t.painter.SetFontSize(float64(t.size))
	tmp := image.NewRGBA(image.Rect(0, 0, 1080, 1920))
	top := 0
	left := 0
	t.painter.Context.SetClip(tmp.Bounds())
	t.painter.Context.SetDst(tmp)
	lines := strings.Split(text, "\n")
	var pt fixed.Point26_6
	for _, line := range lines {
		pt = freetype.Pt(0, top+int(t.painter.Context.PointToFixed(t.painter.FontSize)>>6))
		pt, err = t.painter.Context.DrawString(line, pt)
		if err != nil {
			t.measureError = errors.Wrap(err, "draw string failed")
			return t.Rect
		}
		if int(pt.X)>>6 > left {
			left = int(pt.X) >> 6
		}
		top += int(t.painter.FontSize)
	}
	t.Rect = image.Rect(0, 0, left, top)
	t.background = image.NewRGBA(t.Rect)
	draw.Draw(t.background, t.Rect, tmp, image.ZP, draw.Src)
	return t.Rect
}

func (t *TextView) ViewType() ViewType {
	return ViewTypeTextView
}

// src ttf font file path
func CreateTextView(text string, src string, size int, color string) *TextView {
	tv := &TextView{
		BaseView: &BaseView{
			Attr: Attr{
				AttrKey.Src:       src,
				AttrKey.FontSize:  size,
				AttrKey.FontColor: color,
				AttrKey.Text:      text,
			},
		},
	}
	return tv
}
