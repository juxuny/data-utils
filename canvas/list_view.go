package canvas

import (
	"github.com/pkg/errors"
	"image"
)

type ListView struct {
	*BaseView
	startPoint image.Point
}

func (t *ListView) Draw(img *image.RGBA, vector ...image.Point) error {
	start := t.startPoint
	if len(vector) > 0 {
		start.X += vector[0].X
		start.Y += vector[0].Y
	}
	for i := 0; i < t.ChildNum(); i++ {
		child := t.GetChild(i)
		rect := child.Measure()
		if err := t.GetChild(i).Draw(img, image.Point{
			X: start.X,
			Y: start.Y,
		}); err != nil {
			return errors.Wrapf(err, "draw child(%d) failed", i)
		}
		start.Y += rect.Dy()
	}
	return nil
}

func (t *ListView) Measure() image.Rectangle {
	if t.measured {
		return t.Rect
	}
	height := 0
	width := 0
	for i := 0; i < t.ChildNum(); i++ {
		rect := t.GetChild(i).Measure()
		height += rect.Dy()
		if rect.Dx() > width {
			width = rect.Dx()
		}
	}
	t.Rect = image.Rect(t.startPoint.X, t.startPoint.Y, t.startPoint.X+width, t.startPoint.Y+height)
	t.measured = true
	return t.Rect
}

func (t *ListView) ViewType() ViewType {
	return ViewTypeListView
}

func CreateListView(left, top int, children []View) *ListView {
	lv := &ListView{
		BaseView: &BaseView{
			Children: children,
			Rect:     image.Rectangle{},
		},
	}
	lv.startPoint = image.Point{
		X: left,
		Y: top,
	}
	return lv
}
