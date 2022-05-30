package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ImageWidth      int    `json:"image_width"`
	ImageHeight     int    `json:"image_height"`
	FontSize        int    `json:"font_size"`
	FontColor       string `json:"font_color"`
	BackgroundColor string `json:"background_color"`
	FontName        string `json:"font_name"`
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

	fmt.Printf("config file created: %s\n", filepath)
	return nil
}

func defaultString() (string, error) {
	v, err := json.MarshalIndent(Config{
		ImageWidth:      1280,
		ImageHeight:     720,
		FontSize:        64,
		FontColor:       "#000",
		BackgroundColor: "#fff",
		FontName:        "D2Coding.ttf",
	}, "", "  ")
	if err != nil {
		return "{}", err
	}

	return string(v), nil
}

func Parse(path string) (Config, error) {
	cfg := Config{}
	v, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(v, &cfg)
	return cfg, err
}
