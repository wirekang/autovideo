package logic

import (
	"fmt"
	"math"

	"github.com/samber/lo"
	"github.com/wirekang/autovideo/audio"
	"github.com/wirekang/autovideo/config"
	"github.com/wirekang/autovideo/file"
	"github.com/wirekang/autovideo/image"
	"github.com/wirekang/autovideo/script"
	"github.com/wirekang/autovideo/video"
)

func Start(cfg config.Config, audioDirPath, scriptFilePath, outputFilePath string) error {
	const thumbnailFilePath = "thumb.png"
	const imagesDir = "images"
	const imageFilePrefix = "image_"
	const videoOutputFilePath = "video.mp4"
	const audioOutputFilePath = "audio.mp3"

	file.DeleteOnExit(imagesDir)
	file.DeleteOnExit(videoOutputFilePath)
	file.DeleteOnExit(audioOutputFilePath)

	audioFilePaths, err := file.ListByModTime(audioDirPath)
	if err != nil {
		return fmt.Errorf("can't list audio files: %w", err)
	}

	durations, err := file.MediaLengthSecondsAll(audioFilePaths)
	if err != nil {
		return fmt.Errorf("can't get length of audio files: %w", err)
	}

	thumbnailText, texts, err := script.Parse(scriptFilePath)
	if err != nil {
		return fmt.Errorf("can't parse script files: %w", err)
	}

	if len(texts) != len(durations) {
		return fmt.Errorf("length mismatch: script(%d) <> audio(%d)", len(texts), len(durations))
	}

	err = generateThumbnail(thumbnailText, thumbnailFilePath, cfg)
	if err != nil {
		return fmt.Errorf("can't generate thumbnail image: %w", err)
	}

	err = generateImages(texts, imagesDir, imageFilePrefix, cfg)
	if err != nil {
		return fmt.Errorf("can't generate images: %w", err)
	}

	err = audio.Merge(audioFilePaths, audioOutputFilePath)
	if err != nil {
		return fmt.Errorf("can't merge audios: %w", err)
	}

	err = video.Generate(imagesDir, imageFilePrefix, durations, videoOutputFilePath)
	if err != nil {
		return fmt.Errorf("can't generate video file: %w", err)
	}

	err = video.Merge(videoOutputFilePath, audioOutputFilePath, outputFilePath)
	if err != nil {
		return fmt.Errorf("can't merge audio file with video file: %w", err)
	}

	totalDuration := lo.SumBy(durations, func(v float64) float64 { return v })
	fmt.Println("Total duration:", totalDuration)
	checkDurationDiff(audioOutputFilePath, totalDuration)
	checkDurationDiff(videoOutputFilePath, totalDuration)
	checkDurationDiff(outputFilePath, totalDuration)
	return nil
}

func checkDurationDiff(filePath string, expected float64) {
	audioDuration, err := file.MediaLengthSeconds(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	diff := math.Abs(audioDuration - expected)
	if diff > 0.5 {
		fmt.Printf("file length diff: %vs (%s)\n", diff, filePath)
	}
}

func generateThumbnail(text, outputFilePath string, cfg config.Config) error {
	return image.Generate(text, outputFilePath, image.GenerateOption{
		CanvasWidth:     cfg.ImageWidth,
		CanvasHeight:    cfg.ImageHeight,
		FontFace:        cfg.ThumbnailFontFace,
		FontColor:       cfg.ThumbnailFontColor,
		BackgroundColor: cfg.ThumbnailBackgroundColor,
	})
}

func generateImages(texts []string, outputDirPath, outputPrefix string, cfg config.Config) error {
	err := file.MkdirAllReset(outputDirPath)
	if err != nil {
		return err
	}

	for i, text := range texts {
		err = image.Generate(text, file.Index(outputDirPath, outputPrefix, i, "png"), image.GenerateOption{
			CanvasWidth:     cfg.ImageWidth,
			CanvasHeight:    cfg.ImageHeight,
			FontFace:        cfg.FontFace,
			FontColor:       cfg.FontColor,
			BackgroundColor: cfg.BackgroundColor,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
