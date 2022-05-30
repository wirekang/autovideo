package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/fileutil"
	"github.com/wirekang/autovideo/imageSaver"
	"github.com/wirekang/autovideo/line"
)

func TestImageSaver(t *testing.T) {
	is := assert.New(t)
	CdTempDir()
	outputDir := "images"

	i := imageSaver.New(imageSaver.Option{
		CanvasWidth:     1280,
		CanvasHeight:    720,
		FontPoints:      64,
		FontName:        "D2Coding.ttf",
		FontColor:       "#fff",
		BackgroundColor: "#000",
		OutputDir:       outputDir,
	})

	lines := []line.Line{
		{"Hello", 0},
		{"Word", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
	}
	is.Nil(i.SaveImages(lines))

	count, err := fileutil.Count(outputDir, false)
	is.Nil(err)
	is.Equal(len(lines), count)
}

func TestImageSaverTooLong(t *testing.T) {
	is := assert.New(t)
	CdTempDir()
	imagesOutput := "images"

	i := imageSaver.New(imageSaver.Option{
		CanvasWidth:     1280,
		CanvasHeight:    720,
		FontPoints:      64,
		FontName:        "D2Coding.ttf",
		FontColor:       "#fff",
		BackgroundColor: "#000",
		OutputDir:       imagesOutput,
	})

	lines := []line.Line{
		{"Hello", 0},
		{"Word", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Hello", 0},
		{"Word 한글 테스트 한글 테스트 한글 아주 긴 아주 긴 아주 긴 글자 글자 글자", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
		{"Word 한글 테스트 한글 테스트 한글", 0},
	}
	is.NotNil(i.SaveImages(lines))
}