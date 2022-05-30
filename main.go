package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/imageConcater"
	"github.com/wirekang/autovideo/imageSaver"
	"github.com/wirekang/autovideo/line"
)

const imagesDir = "images"
const imageFilePrefix = "image_"

func main() {
	var configFile string
	var initConfig bool
	var outputFile string
	var linesFile string

	pflag.StringVarP(&configFile, "config", "c", "autovideo.json", "config file path")
	pflag.BoolVar(&initConfig, "init", false, "create default config file")
	pflag.StringVarP(&outputFile, "output", "o", "out.mp4", "output file name")
	pflag.Parse()
	linesFile = pflag.Arg(0)

	if initConfig {
		v := lo.Must(config.Default())
		lo.Must0(os.WriteFile(configFile, []byte(v), 0666))
		fmt.Printf("config file created: %s\n", configFile)
		return
	}

	input, err := os.ReadFile(linesFile)
	if err != nil {
		fmt.Printf("%s <lines.json>\n\n", os.Args[0])
		panic(fmt.Errorf("can't read lines: %w", err))
	}

	lines, err := line.ParseLines(input)
	if err != nil {
		panic(fmt.Errorf("can't parse lines: %w", err))
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		panic(err)
	}

	saver := imageSaver.New(imageSaver.Option{
		OutputDir:       imagesDir,
		ImageFilePrefix: imageFilePrefix,
		CanvasWidth:     cfg.ImageWidth,
		CanvasHeight:    cfg.ImageHeight,
		FontPoints:      float64(cfg.FontSize),
		FontName:        cfg.FontName,
		FontColor:       cfg.FontColor,
		BackgroundColor: cfg.BackgroundColor,
		Lines:           lines,
	})

	err = saver.SaveImages()
	if err != nil {
		panic(fmt.Errorf("can't save images: %w", err))
	}

	concater := imageConcater.New(imageConcater.Option{
		InputDir:        imagesDir,
		ImageFilePrefix: imageFilePrefix,
		OutputFile:      outputFile,
		Lines:           lines,
	})

	err = concater.ConcatImages()
	if err != nil {
		panic(fmt.Errorf("cant' concat images: %w", err))
	}

}
