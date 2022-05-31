package main

import (
	"fmt"

	"github.com/fatih/color"
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

const thumbnailFilePath = "thumb.png"
const imagesDir = "images"
const imageFilePrefix = "image_"
const videoOutputFilePath = "video.mp4"
const audioOutputFilePath = "audio.mp3"

func main() {
	var isInitConfig bool
	var configFilePath string

	var txtFilePath string

	var isReuseScript bool
	var scriptFilePath string

	var outputFilePath string
	var audiosDirPath string

	pflag.BoolVar(&isInitConfig, "init", false, "create default config file")
	pflag.StringVarP(&configFilePath, "config", "c", "autovideo.json", "config file path")

	pflag.StringVar(&txtFilePath, "txt", "txt.txt", "plain text file for script")
	pflag.StringVar(&scriptFilePath, "script", "script.json", "script file")

	pflag.BoolVar(&isReuseScript, "reuse", false, "reuse script file")

	pflag.StringVarP(&outputFilePath, "output", "o", "out.mp4", "output file name")
	pflag.StringVarP(&audiosDirPath, "audios", "a", "audios", "audio files directory")
	pflag.Parse()

	if isInitConfig {
		fmt.Println("create default config file to", configFilePath)
		lo.Must0(config.Init(configFilePath))
		return
	}

	if txtFilePath == "" {
		fmt.Println("no --txt")
		return
	}

	fmt.Println("parse config from", configFilePath)
	cfg := lo.Must(config.Parse(configFilePath))

	ac := audioConcater.New(audioConcater.Option{
		InputDir:   audiosDirPath,
		OutputFile: audioOutputFilePath,
	})

	fmt.Println("merge all audio files")
	lo.Must0(ac.Concat())

	var lines []script.Line
	thumbnail, lineStrings := lo.Must2(script.ExtractThumbnail(txtFilePath))
	if isReuseScript {
		fmt.Println("reuse script from", scriptFilePath)
		lines = lo.Must(script.LoadLines(scriptFilePath))
	} else {
		fmt.Println("run", audioOutputFilePath)
		lo.Must0(fileutil.Run(audioOutputFilePath))

		fmt.Println("generate script from", txtFilePath)
		lines = lo.Must(script.GenerateLines(lineStrings))

		fmt.Println("save script to", scriptFilePath)
		lo.Must0(script.SaveLines(lines, scriptFilePath))
	}

	saver := imageSaver.New(imageSaver.Option{
		CanvasWidth:              cfg.ImageWidth,
		CanvasHeight:             cfg.ImageHeight,
		FontPoints:               float64(cfg.FontSize),
		FontName:                 cfg.FontName,
		FontColor:                cfg.FontColor,
		BackgroundColor:          cfg.BackgroundColor,
		OutputDir:                imagesDir,
		Lines:                    lines,
		ImageFilePrefix:          imageFilePrefix,
		ThumbnailOutputFilePath:  thumbnailFilePath,
		Thumbnail:                thumbnail,
		ThumbnailFontPoints:      float64(cfg.ThumbnailFontSize),
		ThumbnailFontColor:       cfg.ThumbnailFontColor,
		ThumbnailBackgroundColor: cfg.ThumbnailBackgroundColor,
	})
	defer func() {
		fileutil.TryRemove(imagesDir)
	}()

	fmt.Println("generate thumbnail to", thumbnailFilePath)
	lo.Must0(saver.SaveThumbnail())

	fmt.Println("generate images")
	lo.Must0(saver.SaveImages())

	ic := imageConcater.New(imageConcater.Option{
		InputDir:        imagesDir,
		ImageFilePrefix: imageFilePrefix,
		OutputFile:      videoOutputFilePath,
		Lines:           lines,
	})
	fmt.Println("generate video from images")
	lo.Must0(ic.Concat())

	color.Blue(outputFilePath)
	lo.Must0(ffmpegutil.Merge(audioOutputFilePath, videoOutputFilePath, outputFilePath))
}

// todo
// (외부) 전체 플로우 자동화
