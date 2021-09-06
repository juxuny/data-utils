package dict

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"image"
	"io/ioutil"
)

type Painter struct {
	Context  *freetype.Context
	FontSize float64
}

func NewPainter() *Painter {
	p := &Painter{
		Context: freetype.NewContext(),
	}
	return p
}

func (t *Painter) SetFont(fontFile string) error {
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return errors.Wrap(err, "read font file failed")
	}
	//f, err := sfnt.Parse(fontBytes)
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return errors.Wrap(err, "parse font failed")
	}
	t.Context.SetFont(f)
	return nil
}

func (t *Painter) SetFontByFontData(fontBytes []byte) error {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return errors.Wrap(err, "parse font data failed")
	}
	t.Context.SetFont(f)
	return nil
}

func (t *Painter) SetColor(img image.Image) *Painter {
	t.Context.SetSrc(img)
	return t
}

func (t *Painter) SetFontSize(s float64) *Painter {
	t.FontSize = s
	t.Context.SetFontSize(s)
	return t
}
