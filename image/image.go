package image

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type GenerateOption struct {
	CanvasWidth     int
	CanvasHeight    int
	FontFace        font.Face
	FontColor       color.Color
	BackgroundColor color.Color
}

func Generate(text string, outputFilePath string, o GenerateOption) error {
	c := gg.NewContext(o.CanvasWidth, o.CanvasHeight)
	c.SetColor(o.BackgroundColor)
	c.Clear()

	c.SetFontFace(o.FontFace)
	c.SetColor(o.FontColor)
	x := float64(o.CanvasWidth) / 2
	y := float64(o.CanvasHeight) / 2
	ax := 0.5
	ay := 0.5
	w := float64(o.CanvasWidth) * 0.9
	ls := 1.2

	lines := c.WordWrap(text, w)
	for _, line := range lines {
		textW, textH := c.MeasureMultilineString(line, ls)
		if textW > w || textH > float64(o.CanvasHeight)*0.9 {
			return fmt.Errorf("text overflow: %s(%f) > %f", line, textW, w)
		}
	}

	c.DrawStringWrapped(text, x, y, ax, ay, w, ls, gg.AlignCenter)
	return c.SavePNG(outputFilePath)
}
