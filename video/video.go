package video

import (
	"github.com/samber/lo"
	"github.com/wirekang/autovideo/exe"
	"github.com/wirekang/autovideo/file"
)

const ffConcatFilePath = "video.ffconcat"

func Generate(imagesDirPaths, imagesPrefix string, durations []float64, outputFilePath string) error {
	imageFilePaths := lo.Map(durations, func(_ float64, i int) string {
		return file.Index(imagesDirPaths, imagesPrefix, i, "png")
	})

	err := file.WriteFFConcat(imageFilePaths, durations, ffConcatFilePath)
	if err != nil {
		return err
	}

	_, err = exe.Run(
		"ffmpeg",
		"-i",
		ffConcatFilePath,
		"-safe",
		"0",
		"-vcodec",
		"libx264",
		"-an",
		"-y",
		outputFilePath,
	)
	if err != nil {
		return err
	}

	return nil
}

func Merge(videoFilePath, audioFilePath, outputFilePath string) error {
	_, err := exe.Run(
		"ffmpeg",
		"-i",
		videoFilePath,
		"-i",
		audioFilePath,
		"-c:v",
		"copy",
		"-c:a",
		"aac",
		"-y",
		outputFilePath,
	)
	if err != nil {
		return err
	}

	return nil
}
