package ffmpegutil

import (
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wirekang/autovideo/fileutil"
)

func Merge(audio string, video string, output string) error {
	streams := []*ffmpeg_go.Stream{
		ffmpeg_go.Input(audio),
		ffmpeg_go.Input(video),
	}
	fileutil.TryRemove(output)

	defer func() {
		fileutil.TryRemove(audio)
		fileutil.TryRemove(video)
	}()

	return ffmpeg_go.Output(streams, output).Run()
}

func InputOutput(input string, output string) error {
	return ffmpeg_go.Input(input).Output(output).Run()
}
