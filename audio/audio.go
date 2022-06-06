package audio

import (
	"github.com/wirekang/autovideo/exe"
	"github.com/wirekang/autovideo/file"
)

const ffConcatFilePath = "audio.ffconcat"

func Merge(filePaths []string, outputFilePath string) error {
	err := file.WriteFFConcat(filePaths, nil, ffConcatFilePath)
	if err != nil {
		return err
	}

	_, err = exe.Run(
		"ffmpeg",
		"-safe",
		"0",
		"-i",
		ffConcatFilePath,
		"-y",
		outputFilePath,
	)
	if err != nil {
		return err
	}

	return nil
}
