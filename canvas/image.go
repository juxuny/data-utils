package canvas

import (
	"github.com/juxuny/data-utils/lib"
	"golang.org/x/image/font/gofont/goitalic"
)

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

/*
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

*/
