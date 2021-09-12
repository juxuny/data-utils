package canvas

import (
	"github.com/juxuny/data-utils/lib"
	"golang.org/x/image/font/gofont/goitalic"
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
