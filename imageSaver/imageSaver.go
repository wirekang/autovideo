package imageSaver

import (
	"fmt"
	"image/color"
	"os"
	"path"

	"github.com/flopp/go-findfont"
	"github.com/fogleman/gg"
	"github.com/samber/lo"
	"github.com/wirekang/autovideo/colorutil"
	"github.com/wirekang/autovideo/line"
	"golang.org/x/image/font"
)

type Interface interface {
	SaveImages(lines []line.Line) error
}

type inner struct {
	canvasW, canvasH int
	fontColor        color.Color
	bgColor          color.Color
	fontFace         font.Face
	outputDir        string
	filePrefix       string
}

type Option struct {
	CanvasWidth,
	CanvasHeight int
	FontPoints float64
	FontName   string
	FontColor,
	BackgroundColor string
	OutputDir string
}

func New(o Option) Interface {
	fontColor := lo.Must(colorutil.ParseHexColor(o.FontColor))
	bgColor := lo.Must(colorutil.ParseHexColor(o.BackgroundColor))
	fontFace := lo.Must(loadFontFace(o.FontName, o.FontPoints))
	lo.Must0(os.MkdirAll(o.OutputDir, os.ModeDir))

	return &inner{
		outputDir:  o.OutputDir,
		canvasW:    o.CanvasWidth,
		canvasH:    o.CanvasHeight,
		fontColor:  fontColor,
		bgColor:    bgColor,
		fontFace:   fontFace,
		filePrefix: "output",
	}
}

func (i *inner) SaveImages(lines []line.Line) error {
	c := gg.NewContext(i.canvasW, i.canvasH)
	c.SetFontFace(i.fontFace)

	for j, line := range lines {
		c.SetColor(i.bgColor)
		c.Clear()
		c.SetColor(i.fontColor)
		err := draw(c, line.Text)
		if err != nil {
			return fmt.Errorf("can't draw line %d(%s): %w", j, fmtLine(line), err)
		}

		err = save(c, i.outputDir, i.filePrefix, j)
		if err != nil {
			return fmt.Errorf("can't save line %d(%s): %w", j, fmtLine(line), err)
		}
	}
	return nil
}

func fmtLine(l line.Line) string {
	return fmt.Sprintf("%.5s...", l.Text)
}

func save(c *gg.Context, outputDir string, filePrefix string, i int) error {
	dst := path.Join(outputDir, fmt.Sprintf("%s%d.png", filePrefix, i))
	return c.SavePNG(dst)
}

func draw(c *gg.Context, text string) error {
	canvasW := c.Width()
	canvasH := c.Height()
	textW, textH := c.MeasureString(text)
	textX := (float64(canvasW) - (textW)) / 2
	textY := (float64(canvasH) - (textH)) / 2
	if textW >= float64(canvasW)*0.9 {
		return fmt.Errorf("text is too long")
	}

	c.DrawString(text, textX, textY)
	return nil
}

func loadFontFace(name string, points float64) (font.Face, error) {
	p, err := findfont.Find(name)
	if err != nil {
		return nil, fmt.Errorf("can't find %s: %w", name, err)
	}

	face, err := gg.LoadFontFace(p, points)
	if err != nil {
		return nil, fmt.Errorf("can't load font %s: %w", p, err)
	}

	return face, nil
}
