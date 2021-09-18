package canvas

import (
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"image"
)

type CenterType uint8

const (
	CenterTypeAll        = CenterType(0)
	CenterTypeHorizontal = CenterType(1)
	CenterTypeVertical   = CenterType(2)
)

type CenterLayout struct {
	*BaseView
	measureError error
}

func (t *CenterLayout) Draw(img *image.RGBA, vector ...image.Point) error {
	if t.measureError != nil {
		return t.measureError
	}
	start := image.Point{
		X: t.Rect.Min.X,
		Y: t.Rect.Min.Y,
	}
	if len(vector) > 0 {
		start.X += vector[0].X
		start.Y += vector[0].Y
	}
	centerType, ok := t.Attr[AttrKey.CenterType].(CenterType)
	if !ok {
		return invalidAttr(AttrKey.CenterType, t.Attr[AttrKey.CenterType])
	}
	for i, c := range t.Children {
		rect := c.Measure()
		log.Debug(rect.Min)
		vector := image.Point{}.Add(start)
		vector = vector.Sub(rect.Min)
		if centerType == CenterTypeAll || centerType == CenterTypeHorizontal {
			vector.X += t.Rect.Dx() >> 1
			vector.X -= rect.Dx() >> 1
			vector.Y += rect.Min.Y
		}
		if centerType == CenterTypeAll || centerType == CenterTypeVertical {
			vector.Y += t.Rect.Dy() >> 1
			vector.Y -= rect.Dy() >> 1
			vector.X += rect.Min.X
		}
		if err := c.Draw(img, vector); err != nil {
			return errors.Wrapf(err, "draw child(%d) failed", i)
		}
	}
	return nil
}

func (t *CenterLayout) Measure() image.Rectangle {
	if t.measured {
		return t.Rect
	}
	for _, v := range t.Children {
		v.Measure()
	}
	t.measured = true
	return t.Rect
}

func CreateCenterLayout(centerType CenterType, start image.Point, width int, height int, child View) *CenterLayout {
	cl := &CenterLayout{
		BaseView: &BaseView{
			Rect: image.Rect(start.X, start.Y, start.X+width, start.Y+height),
			Attr: Attr{
				AttrKey.CenterType: centerType,
			},
			Children: []View{
				child,
			},
		},
	}
	return cl
}

func (t *CenterLayout) ViewType() ViewType {
	return ViewTypeCenterLayout
}
