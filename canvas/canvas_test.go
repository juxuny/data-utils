package canvas

import (
	"image"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	canvas := NewCanvas(1920, 810*2)
	img := CreateImageView("tmp/cover.jpeg", 1920, 810, ImageTypeJpeg)
	if err := canvas.Draw(img); err != nil {
		t.Fatal(err)
	}

	lv := CreateListView(0, 0, []View{
		CreateImageView("tmp/cover.jpeg", 1920, 810, ImageTypeJpeg),
		//CreateImageView("tmp/cover.jpeg", 1920, 810, ImageTypeJpeg),
		CreateCenterLayout(CenterTypeAll, image.Point{X: 0, Y: 0}, 1920, 810,
			CreateBox(image.Rect(0, 0, 1920, 810),
				CreateListView(0, 0, []View{
					CreateTextView("Hello\nWorld !!!\nMy name is Juxuny Wu", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#FFFFFF"),
					CreateTextView("Daily Help\nWorld !!!\nMy name is Juxuny Wu", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#FFFFFF"),
					CreateTextView("Hello\nWorld !!!\nMy name is Juxuny Wu", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#00FFFF"),
					CreateTextView("Hello\nWorld !!!\nMy name is Juxuny Wu", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#FFFFFF"),
					CreateTextView("Hello\nWorld !!!", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#FFFFFF"),
					CreateTextView("Hello\nWorld !!!", "tmp/No.73ShangShouFenBiTi-2.ttf", 42, "#FFFFFF"),
				}),
			),
		),
	})

	if err := canvas.Draw(lv); err != nil {
		t.Fatal(err)
	}

	if err := canvas.Save("tmp/hello.jpg", ImageTypeJpeg); err != nil {
		t.Fatal(err)
	}
}
