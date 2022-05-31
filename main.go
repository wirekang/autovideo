package main

import (
	"fmt"
	"os/exec"

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
	var isInitConfig bool
	var configFilePath string

	var txtFilePath string

	var outputFilePath string
	var audiosDirPath string

	pflag.BoolVar(&isInitConfig, "init", false, "create default config file")
	pflag.StringVarP(&configFilePath, "config", "c", "autovideo.json", "config file path")

	pflag.StringVar(&txtFilePath, "txt", "", "plain text file for script")

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

	lo.Must0(exec.Command("cmd", "/c", "start", audioOutputFilePath).Run())
	fmt.Println("generate script from", txtFilePath)
	lines := lo.Must(script.Generate(txtFilePath))

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

	fmt.Println("merge video with audio to", outputFilePath)
	lo.Must0(ffmpegutil.Merge(audioOutputFilePath, videoOutputFilePath, outputFilePath))
}
