package dict

import (
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
		CreateTextView("Hello\nWorld !!!", "tmp/Cronos-Pro-Bold_12435.ttf", 32, "#FFFFFF"),
		CreateTextView("Hello\nWorld !!!", "tmp/Cronos-Pro-Bold_12435.ttf", 32, "#FFFFFF"),
		CreateTextView("Hello\nWorld !!!", "tmp/Cronos-Pro-Bold_12435.ttf", 32, "#FFFFFF"),
	})

	if err := canvas.Draw(lv); err != nil {
		t.Fatal(err)
	}

	if err := canvas.Save("tmp/hello.jpg", ImageTypeJpeg); err != nil {
		t.Fatal(err)
	}
}
