package dict

import (
	"bufio"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

type Canvas struct {
	Background *image.RGBA
	Painter    *Painter
}

func NewCanvas(width, height int) (c *Canvas) {
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), image.Transparent, image.ZP, draw.Src)
	return &Canvas{
		Background: background,
	}
}

func (t *Canvas) SetBackgroundColor(img image.Image) {
	draw.Draw(t.Background, t.Background.Bounds(), img, image.ZP, draw.Src)
}

func (t *Canvas) DrawColor(c string) error {
	s, err := convertHexToRGBA(c)
	if err != nil {
		return errors.Wrap(err, "invalid color value: "+c)
	}
	t.SetBackgroundColor(image.NewUniform(s))
	return nil
}

func (t *Canvas) DrawImageFromFile(point image.Point, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "draw image from file failed, can't open the file "+file)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return errors.Wrap(err, "invalid image file")
	}
	draw.Draw(t.Background, image.Rect(point.X, point.Y, t.Background.Rect.Max.X, t.Background.Rect.Max.Y), img, image.ZP, draw.Src)
	return nil
}

func (t *Canvas) SetPainter(p *Painter) *Canvas {
	t.Painter = p
	return t
}

func (t *Canvas) DrawText(text string, left, top int) error {
	lines := strings.Split(text, "\n")
	var pt fixed.Point26_6
	var err error
	for _, line := range lines {
		t.Painter.Context.SetClip(t.Background.Bounds())
		t.Painter.Context.SetDst(t.Background)
		pt = freetype.Pt(left, top+int(t.Painter.Context.PointToFixed(t.Painter.FontSize)>>6))
		pt, err = t.Painter.Context.DrawString(line, pt)
		if err != nil {
			return err
		}
		top += int(t.Painter.FontSize)
	}
	return nil
}

func (t *Canvas) Save(file string, imageType ...ImageType) error {
	out, err := os.Create(file)
	if err != nil {
		return errors.Wrap(err, "save failed")
	}
	defer func() {
		_ = out.Close()
	}()
	buf := bufio.NewWriter(out)

	it := ImageTypeJpeg
	if len(imageType) > 0 {
		it = imageType[0]
	}

	if it == ImageTypeJpeg {
		if err := jpeg.Encode(buf, t.Background, &jpeg.Options{Quality: 90}); err != nil {
			return errors.Wrap(err, "save failed")
		}
	} else {
		if err := png.Encode(buf, t.Background); err != nil {
			return errors.Wrap(err, "save failed")
		}
	}

	if err := buf.Flush(); err != nil {
		return errors.Wrap(err, "flush failed")
	}
	return nil

}

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

type ImageType uint8

const (
	ImageTypePng  = ImageType(1)
	ImageTypeJpeg = ImageType(2)
)

type Options struct {
	FontFace   string
	FontSize   float64
	FontColor  string
	FontFile   string
	FontBytes  []byte
	Width      int
	Height     int
	ImageType  ImageType
	Padding    Padding
	Background string
}

type Padding struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func defaultOptions() *Options {
	return &Options{
		FontSize:  12,
		FontColor: "#000000",
		FontBytes: goitalic.TTF,
		FontFile:  "",
		Width:     300,
		Height:    300,
		ImageType: ImageTypeJpeg,
		Padding:   Padding{},
	}
}

func mergeOptions(opts []*Options) *Options {
	ret := defaultOptions()
	for _, o := range opts {
		if o.ImageType == 0 {
			o.ImageType = ImageTypeJpeg
		}
		if o.FontSize == 0 {
			o.FontSize = 12
		}
		if o.FontFace != ret.FontFace {
			ret.FontFace = o.FontFace
		}
		if ret.FontSize != o.FontSize {
			ret.FontSize = o.FontSize
		}
		if ret.FontColor != o.FontColor {
			ret.FontColor = o.FontColor
		}
		if !lib.Byte.CompareSlice(ret.FontBytes, o.FontBytes) {
			ret.FontBytes = o.FontBytes
		}
		if ret.FontFile != o.FontFile {
			ret.FontFile = o.FontFile
		}
		if ret.Width != o.Width {
			ret.Width = o.Width
		}
		if ret.Height != o.Height {
			ret.Height = o.Height
		}
		if ret.ImageType != o.ImageType {
			ret.ImageType = o.ImageType
		}
		if ret.Padding != o.Padding {
			ret.Padding = o.Padding
		}
		if ret.Background != o.Background {
			ret.Background = o.Background
		}
	}
	return ret
}

func GenerateWordImage(outFile string, title string, words []Word, opts ...*Options) (err error) {
	option := mergeOptions(opts)
	c := NewCanvas(option.Width, option.Height)
	p := NewPainter()
	setFont := false
	if option.FontFile != "" {
		if err := p.SetFont(option.FontFile); err != nil {
			log.Warn(err)
		} else {
			setFont = true
		}
	}

	if strings.Index(option.Background, "#") == 0 {
		// set background color
		if err := c.DrawColor(option.Background); err != nil {
			return errors.Wrap(err, "init background color failed")
		}
	} else {
		// draw background image
		if err := c.DrawImageFromFile(image.Point{
			X: 0,
			Y: 0,
		}, option.Background); err != nil {
			return errors.Wrap(err, "generate cover failed")
		}
	}

	// init font face
	if !setFont && option.FontBytes != nil {
		if err := p.SetFontByFontData(option.FontBytes); err != nil {
			return errors.Wrap(err, "set font failed")
		} else {
			setFont = true
		}
	}
	c.SetPainter(p)
	fontColor, err := convertHexToRGBA(option.FontColor)
	if err != nil {
		return errors.Wrap(err, "invalid color format: "+option.FontColor)
	}
	p.SetColor(image.NewUniform(fontColor))
	p.SetFontSize(option.FontSize)

	text := ""
	if title != "" {
		text += title + "\n"
	}
	for _, w := range words {
		text += w.Name + "\n"
	}
	if err := c.DrawText(text, option.Padding.Left, option.Padding.Top); err != nil {
		return errors.Wrap(err, "generate image failed")
	}

	if option.ImageType == ImageTypeJpeg {
		return c.Save(outFile)
	}
	return nil
}
