package canvas

import (
	"fmt"
	"image"
)

var DPI = float64(72)

type ViewType string

const (
	ViewTypeImageView    = ViewType("image-view")
	ViewTypeListView     = ViewType("list-view")
	ViewTypeTextView     = ViewType("text-view")
	ViewTypeBox          = ViewType("box")
	ViewTypeCenterLayout = ViewType("center-layout")
)

type View interface {
	Measured() bool
	Draw(img *image.RGBA, vector ...image.Point) error
	Measure() image.Rectangle
	GetChild(index int) View
	ChildNum() int
	Width() int
	Height() int
	GetAttr() Attr
	ViewType() ViewType
	AddChild(index int, child View) error
	AppendChild(child View) error
}

type BaseView struct {
	Children []View
	Rect     image.Rectangle
	Attr     Attr
	measured bool
}

func (t *BaseView) Width() int {
	return t.Rect.Dx()
}

func (t *BaseView) Height() int {
	return t.Rect.Dy()
}

func (t *BaseView) ChildNum() int {
	return len(t.Children)
}

func (t *BaseView) GetChild(index int) View {
	return t.Children[index]
}

func (t *BaseView) GetAttr() Attr {
	return t.Attr
}

func (t *BaseView) SetAttr(key string, value interface{}) {
	t.Attr[key] = value
}

func (t *BaseView) AppendChild(child View) error {
	t.Children = append(t.Children, child)
	return nil
}

func (t *BaseView) AddChild(index int, child View) error {
	newChildren := make([]View, 0)
	for i := 0; i < len(t.Children); i++ {
		if i == index {
			newChildren = append(newChildren, child, t.Children[i])
			continue
		}
		newChildren = append(newChildren, t.Children[i])
	}
	if len(newChildren) == len(t.Children) {
		return fmt.Errorf("invalid index: %d", index)
	}
	t.Children = newChildren
	return nil
}

func (t *BaseView) Measured() bool {
	return t.measured
}
