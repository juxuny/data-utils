package dict

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	c := NewCanvas(500, 300)
	c.SetBackgroundColor(image.NewUniform(color.Black))
	p := NewPainter()
	if err := p.SetFont("tmp/Cronos-Pro-Bold_12435.ttf"); err != nil {
		t.Fatal(err)
	}
	c.SetPainter(p)
	p.SetColor(image.NewUniform(color.RGBA{R: 0xff, G: 0xf0, B: 0xcf, A: 0xff}))
	p.SetFontSize(16)
	if err := c.DrawText("Hello\nWorld", 0, 0); err != nil {
		t.Fatal(err)
	}
	if err := c.Save("tmp/hello.png"); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateWordImage(t *testing.T) {
	if err := GenerateWordImage("tmp/word.jpg", "本期重点单词", []Word{
		{Name: "word"},
		{Name: "word"},
		{Name: "word"},
	}, &Options{
		FontSize:  32,
		FontColor: "#FFF0CF",
		FontFile:  "tmp/No.73ShangShouFenBiTi-2.ttf",
		Width:     1920,
		Height:    810,
		ImageType: ImageTypeJpeg,
		Padding:   Padding{Top: int(math.Trunc(810 * 0.10)), Left: int(math.Trunc(1920 * 0.1))},
	}); err != nil {
		t.Fatal(err)
	}
}
