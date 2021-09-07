package canvas

import (
	"github.com/pkg/errors"
	"image"
	"image/draw"
)

type Box struct {
	*BaseView
	measureError error
}

func (t *Box) ViewType() ViewType {
	return ViewTypeBox
}

func CreateBox(rect image.Rectangle, child View) *Box {
	b := &Box{
		BaseView: &BaseView{
			Children: []View{
				child,
			},
			Rect:     rect,
			Attr:     Attr{},
			measured: false,
		},
	}
	return b
}

func (t *Box) Measure() image.Rectangle {
	if t.measured {
		return t.Rect
	}
	for _, c := range t.Children {
		c.Measure()
	}
	t.measured = true
	return t.Rect
}

func (t *Box) Draw(img *image.RGBA, vector ...image.Point) error {
	if t.measureError != nil {
		return t.measureError
	}
	start := t.Rect.Min
	if len(vector) > 0 {
		start.X += vector[0].X
		start.Y += vector[0].Y
	}
	rect := image.Rect(start.X, start.Y, start.X+t.Rect.Dx(), start.Y+t.Rect.Dy())
	tmp := image.NewRGBA(image.Rect(0, 0, t.Rect.Max.X, t.Rect.Max.Y))
	for _, v := range t.Children {
		if !v.Measured() {
			v.Measure()
		}
		if err := v.Draw(tmp); err != nil {
			return errors.Wrap(err, "draw children failed")
		}
	}
	draw.Draw(img, rect, tmp, image.Point{}, draw.Over)
	return nil
}
