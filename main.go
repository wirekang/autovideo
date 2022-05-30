package main

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/wirekang/autovideo/audioConcater"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/ffmpegutil"
	"github.com/wirekang/autovideo/fileutil"
	"github.com/wirekang/autovideo/imageConcater"
	"github.com/wirekang/autovideo/imageSaver"
	"github.com/wirekang/autovideo/script"
)

const imagesDir = "images"
const imageFilePrefix = "image_"
const videoOutputFilePath = "video.mp4"
const audioOutputFilePath = "audio.mp3"

func main() {
	var configFilePath string
	var isInitConfig bool

	var outputFilePath string
	var scriptFileDir string
	var audiosDirPath string

	pflag.StringVarP(&configFilePath, "config", "c", "autovideo.json", "config file path")
	pflag.BoolVar(&isInitConfig, "init", false, "create default config file")
	pflag.StringVarP(&outputFilePath, "output", "o", "out.mp4", "output file name")
	pflag.StringVarP(&audiosDirPath, "audios", "a", "audios", "audio files directory")
	pflag.StringVarP(&scriptFileDir, "script", "s", "script.json", "script file")
	pflag.Parse()

	if isInitConfig {
		lo.Must0(config.Init(configFilePath))
		return
	}

	lines, err := script.Parse(scriptFileDir)
	if err != nil {
		panic(fmt.Errorf("can't parse script file: %w", err))
	}

	cfg, err := config.Parse(configFilePath)
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
	defer func() {
		fileutil.TryRemove(imagesDir)
	}()

	err = saver.SaveImages()
	if err != nil {
		panic(fmt.Errorf("can't save images: %w", err))
	}

	ic := imageConcater.New(imageConcater.Option{
		InputDir:        imagesDir,
		ImageFilePrefix: imageFilePrefix,
		OutputFile:      videoOutputFilePath,
		Lines:           lines,
	})

	err = ic.Concat()
	if err != nil {
		panic(fmt.Errorf("cant' concat images: %w", err))
	}

	ac := audioConcater.New(audioConcater.Option{
		InputDir:   audiosDirPath,
		OutputFile: audioOutputFilePath,
	})

	err = ac.Concat()
	if err != nil {
		panic(fmt.Errorf("cant' concat audios: %w", err))
	}

	err = ffmpegutil.Merge(audioOutputFilePath, videoOutputFilePath, outputFilePath)
	if err != nil {
		panic(err)
	}

}
