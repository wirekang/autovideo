package config

import (
	"encoding/json"
	"fmt"
	icolor "image/color"
	"os"

	"github.com/flopp/go-findfont"
	"github.com/fogleman/gg"
	"github.com/wirekang/autovideo/color"
	"golang.org/x/image/font"
)

type configJSON struct {
	ImageWidth               int    `json:"image_width"`
	ImageHeight              int    `json:"image_height"`
	FontSize                 int    `json:"font_size"`
	FontColor                string `json:"font_color"`
	BackgroundColor          string `json:"background_color"`
	FontName                 string `json:"font_name"`
	ThumbnailFontSize        int    `json:"thumbnail_font_size"`
	ThumbnailFontColor       string `json:"thumbnail_font_color"`
	ThumbnailBackgroundColor string `json:"thumbnail_background_color"`
}

type Config struct {
	ImageWidth  int
	ImageHeight int

	FontColor       icolor.Color
	BackgroundColor icolor.Color
	FontFace        font.Face

	ThumbnailFontColor       icolor.Color
	ThumbnailBackgroundColor icolor.Color
	ThumbnailFontFace        font.Face
}

func Init(filepath string) (err error) {
	v, err := defaultString()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, []byte(v), 0666)
	if err != nil {
		return err
	}

	return nil
}

func defaultString() (string, error) {
	v, err := json.MarshalIndent(configJSON{
		ImageWidth:               1280,
		ImageHeight:              720,
		FontSize:                 64,
		FontColor:                "#fff",
		BackgroundColor:          "#222",
		FontName:                 "D2Coding.ttf",
		ThumbnailFontSize:        128,
		ThumbnailFontColor:       "#fff",
		ThumbnailBackgroundColor: "#222",
	}, "", "  ")
	if err != nil {
		return "{}", err
	}

	return string(v), nil
}

func Parse(path string) (Config, error) {
	v, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	j := configJSON{}
	err = json.Unmarshal(v, &j)
	if err != nil {
		return Config{}, err
	}

	return convert(j)
}

func convert(j configJSON) (Config, error) {
	fontColor, err := color.ParseHexColor(j.FontColor)
	if err != nil {
		return Config{}, err
	}

	bgColor, err := color.ParseHexColor(j.BackgroundColor)
	if err != nil {
		return Config{}, err
	}

	tFontColor, err := color.ParseHexColor(j.ThumbnailFontColor)
	if err != nil {
		return Config{}, err
	}

	tBgColor, err := color.ParseHexColor(j.ThumbnailBackgroundColor)
	if err != nil {
		return Config{}, err
	}

	fontFace, err := loadFontFace(j.FontName, float64(j.FontSize))
	if err != nil {
		return Config{}, err
	}

	tFontFace, err := loadFontFace(j.FontName, float64(j.ThumbnailFontSize))

	return Config{
		ImageWidth:  j.ImageWidth,
		ImageHeight: j.ImageHeight,

		FontColor:       fontColor,
		BackgroundColor: bgColor,
		FontFace:        fontFace,

		ThumbnailFontColor:       tFontColor,
		ThumbnailBackgroundColor: tBgColor,
		ThumbnailFontFace:        tFontFace,
	}, nil
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
