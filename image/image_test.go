package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/file"
)

func TestGenerate(t *testing.T) {
	is := assert.New(t)

	cfg, err := config.Parse("../sample/config.json")
	is.Nil(err)

	const output = "output.png"
	err = Generate("This is very very very long long long text for github.com/wirekang/autovideoasdfdf", output, GenerateOption{
		CanvasWidth:     cfg.ImageWidth,
		CanvasHeight:    cfg.ImageHeight,
		FontFace:        cfg.FontFace,
		FontColor:       cfg.FontColor,
		BackgroundColor: cfg.BackgroundColor,
	})
	is.Nil(err)
	is.FileExists(output)
	file.DeleteOnExit(output)
	file.Delete()
}
