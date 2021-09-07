package dict

import (
	"bufio"
	"github.com/golang/freetype"
	"github.com/pkg/errors"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type Canvas struct {
	Background *image.RGBA
	Painter    *Painter
}

func NewCanvas(width, height int) (c *Canvas) {
	rect := image.Rect(0, 0, width, height)
	background := image.NewRGBA(rect)
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

func (t *Canvas) Draw(v View) error {
	if !v.Measured() {
		v.Measure()
	}
	return v.Draw(t.Background)
}
