package canvas

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/juxuny/data-utils/lib"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"strings"
	"unicode"
)

type TextView struct {
	*BaseView
	measureError error
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
	draw.Draw(img, image.Rect(start.X, start.Y, start.X+t.Rect.Dx(), start.Y+t.Rect.Dy()), t.background, image.ZP, draw.Over)
	return nil
}

func (t *TextView) getTextHeight() int {
	//scale := float64(t.size) / float64(t.painter.Font.FUnitsPerEm())
	//bounds := t.painter.Font.Bounds(fixed.Int26_6(t.painter.Font.FUnitsPerEm()))
	//height := int(float64(bounds.Max.Y-bounds.Min.Y) * scale)
	fc := truetype.NewFace(t.painter.Font, &truetype.Options{
		Size:    t.painter.FontSize,
		DPI:     DPI,
		Hinting: font.HintingNone,
	})
	m := fc.Metrics()
	//log.Debug(m.Descent.Ceil(), m.Ascent.Ceil(), m.Height.Ceil(), t.size)
	return m.Height.Ceil() + m.Descent.Ceil()
	//return height
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
	t.painter = NewPainter()
	size, isInt := t.Attr[AttrKey.FontSize].(int)
	if !isInt {
		t.measureError = fmt.Errorf("%s is not an int", AttrKey.FontSize)
		return t.Rect
	}
	t.painter.SetFontSize(float64(size))
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
	wordWrap, _ := t.Attr[AttrKey.WordWrap].(bool)
	//tmp := image.NewRGBA(image.Rect(0, 0, 1080, 1920))
	top := 0
	left := 0
	lines := strings.Split(text, "\n")
	var pt fixed.Point26_6
	d := &font.Drawer{
		Face: truetype.NewFace(t.painter.Font, &truetype.Options{
			Size:    t.painter.FontSize,
			DPI:     DPI,
			Hinting: font.HintingNone,
		}),
	}
	//measure bound
	textHeight := t.getTextHeight()
	if wordWrap {
		fixedWidth := t.Rect.Dx()
		t.Rect = image.Rect(0, 0, fixedWidth, 0)
		for _, line := range lines {
			wrapLine := t.wrapLine(line)
			t.Rect.Max.Y = len(wrapLine) * textHeight
		}
	} else {
		t.Rect = image.Rect(0, 0, 0, 0)
		for _, line := range lines {
			_, result := d.BoundString(line)
			width := result.Ceil()
			if width > t.Rect.Max.X {
				t.Rect.Max.X = width
			}
			t.Rect.Max.Y += textHeight
		}
	}
	//log.Debug(t.Rect.Dx(), t.Rect.Dy())
	t.background = image.NewRGBA(t.Rect)
	t.painter.Context.SetClip(t.background.Bounds())
	t.painter.Context.SetDst(t.background)
	for _, line := range lines {
		var wrapLine []string
		if wordWrap {
			wrapLine = t.wrapLine(line)
		} else {
			wrapLine = []string{line}
		}
		for _, wl := range wrapLine {
			pt = freetype.Pt(0, top+int(t.painter.Context.PointToFixed(t.painter.FontSize)>>6))
			pt, err = t.painter.Context.DrawString(wl, pt)
			if err != nil {
				t.measureError = errors.Wrap(err, "draw string failed")
				return t.Rect
			}
			if int(pt.X)>>6 > left {
				left = int(pt.X) >> 6
			}
			top += textHeight
		}
	}
	return t.Rect
}

// 给一行文本自动换行
func (t *TextView) wrapLine(line string) []string {
	l := lib.String.SplitWithStopFunc(line, unicode.IsLetter)
	fixedWidth := t.Width()
	i := 1
	ret := make([]string, 0)
	buf := l[0]
	for i < len(l) {
		w := t.painter.MeasureTextWidth(buf + l[i])
		if fixedWidth < w {
			ret = append(ret, buf)
			buf = l[i]
			i++
			continue
		}
		buf += l[i]
		i++
	}
	if len(buf) > 0 {
		ret = append(ret, buf)
	}
	return ret
}

func (t *TextView) ViewType() ViewType {
	return ViewTypeTextView
}

// src ttf font file path
func CreateTextView(text string, fontFileSrc string, size int, color string) *TextView {
	tv := &TextView{
		BaseView: &BaseView{
			Attr: Attr{
				AttrKey.Src:       fontFileSrc,
				AttrKey.FontSize:  size,
				AttrKey.FontColor: color,
				AttrKey.Text:      text,
			},
		},
	}
	return tv
}

// 创建自动换行的TextView
func CreateWrapTextView(text string, fontFileSrc string, size int, color string, width int, attr *Attr) *TextView {

	tv := &TextView{
		BaseView: &BaseView{
			Attr: Attr{
				AttrKey.Src:       fontFileSrc,
				AttrKey.FontSize:  size,
				AttrKey.FontColor: color,
				AttrKey.Text:      text,
				AttrKey.WordWrap:  true,
			},
			Rect: image.Rect(0, 0, width, 0),
		},
	}
	if attr != nil {
		for k, v := range *attr {
			tv.Attr[k] = v
		}
	}
	return tv
}
