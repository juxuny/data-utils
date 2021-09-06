package dict

import (
	"bufio"
	"github.com/pkg/errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageView struct {
	*BaseView
	measureError error
	img          image.Image
}

func (t *ImageView) Draw(img *image.RGBA, vector ...image.Point) error {
	if t.measureError != nil {
		return t.measureError
	}
	start := t.Rect.Min
	if len(vector) > 0 {
		start.X += vector[0].X
		start.Y += vector[0].Y
	}
	draw.Draw(img, image.Rect(start.X, start.Y, start.X+t.Rect.Dx(), start.Y+t.Rect.Dy()), t.img, image.ZP, draw.Src)
	return nil
}

func (t *ImageView) Measure() image.Rectangle {
	if t.measured {
		return t.Rect
	}
	src, isString := t.Attr[AttrKey.Src].(string)
	if !isString {
		t.measureError = invalidAttr(AttrKey.Src, src)
		return t.Rect
	}
	imageType, isImageType := t.Attr[AttrKey.ImageType].(ImageType)
	if !isImageType {
		t.measureError = invalidAttr(AttrKey.ImageType, imageType)
	}
	f, err := os.Open(src)
	if err != nil {
		t.measureError = errors.Wrap(err, "read image file failed")
		return t.Rect
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	if imageType == ImageTypeJpeg {
		img, err := jpeg.Decode(buf)
		if err != nil {
			t.measureError = errors.Wrap(err, "parse jpeg image failed")
			return t.Rect
		}
		t.img = img
	} else if imageType == ImageTypePng {
		img, err := png.Decode(buf)
		if err != nil {
			t.measureError = errors.Wrap(err, "parse png image failed")
			return t.Rect
		}
		t.img = img
	} else {
		t.measureError = invalidAttr(AttrKey.ImageType, imageType)
		return t.Rect
	}
	t.Rect = t.img.Bounds()
	t.measured = true
	return t.Rect
}

func CreateImageView(src string, width, height int, imageType ImageType) *ImageView {
	ret := &ImageView{
		BaseView: &BaseView{
			Rect: image.Rect(0, 0, width, height),
			Attr: Attr{
				AttrKey.Src:       src,
				AttrKey.ImageType: imageType,
			},
		},
	}
	return ret
}

func (t *ImageView) ViewType() ViewType {
	return ViewTypeImageView
}
