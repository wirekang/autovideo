package main

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/file"
	"github.com/wirekang/autovideo/logic"
)

func main() {
	var isInitConfig bool
	var configFilePath string

	var scriptFilePath string

	var outputFilePath string
	var audiosDirPath string

	var isDebug bool

	pflag.BoolVar(&isInitConfig, "init", false, "create default config file")
	pflag.StringVarP(&configFilePath, "config", "c", "autovideo.json", "config file path")

	pflag.StringVar(&scriptFilePath, "script", "script.txt", "script file")

	pflag.StringVarP(&outputFilePath, "output", "o", "out.mp4", "output file name")
	pflag.StringVarP(&audiosDirPath, "audios", "a", "audios", "audio files directory")

	pflag.BoolVar(&isDebug, "debug", false, "debug mode")
	pflag.Parse()

	if isInitConfig {
		fmt.Println("create default config file to", configFilePath)
		lo.Must0(config.Init(configFilePath))
		return
	}

	file.DoNotDeleteTempFiles = isDebug

	file.ListenDeleteOnExit()
	defer file.Delete()

	cfg := lo.Must(config.Parse(configFilePath))
	lo.Must0(logic.Start(cfg, audiosDirPath, scriptFilePath, outputFilePath))
}
