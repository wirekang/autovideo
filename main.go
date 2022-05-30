package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/imageSaver"
	"github.com/wirekang/autovideo/line"
)

const imagesDir = "images"

func main() {
	var configPath string
	var initConfig bool
	var outputPath string
	var linesJsonPath string

	pflag.StringVarP(&configPath, "config", "c", "autovideo.json", "config file path")
	pflag.BoolVar(&initConfig, "init", false, "create default config file")
	pflag.StringVarP(&outputPath, "output", "o", "out.mp4", "output file name")
	pflag.Parse()
	linesJsonPath = pflag.Arg(0)

	if initConfig {
		v := lo.Must(config.Default())
		lo.Must0(os.WriteFile(configPath, []byte(v), 0666))
		fmt.Printf("config file created: %s\n", configPath)
		return
	}

	input, err := os.ReadFile(linesJsonPath)
	if err != nil {
		fmt.Printf("%s <lines.json>\n\n", os.Args[0])
		panic(fmt.Errorf("can't read lines: %w", err))
	}

	lines, err := line.ParseLines(input)
	if err != nil {
		panic(fmt.Errorf("can't parse lines: %w", err))
	}

	cfg, err := config.Parse(configPath)
	if err != nil {
		panic(err)
	}

	i := imageSaver.New(imageSaver.Option{
		CanvasWidth:     cfg.ImageWidth,
		CanvasHeight:    cfg.ImageHeight,
		FontPoints:      float64(cfg.FontSize),
		FontName:        cfg.FontName,
		FontColor:       cfg.FontColor,
		BackgroundColor: cfg.BackgroundColor,
		OutputDir:       imagesDir,
	})

	err = i.SaveImages(lines)
	if err != nil {
		panic(fmt.Errorf("can't save images: %w", err))
	}

}
