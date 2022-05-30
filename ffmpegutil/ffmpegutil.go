package ffmpegutil

import ffmpeg_go "github.com/u2takey/ffmpeg-go"

func Merge(audio string, video string, output string) error {
	streams := []*ffmpeg_go.Stream{
		ffmpeg_go.Input(audio),
		ffmpeg_go.Input(video),
	}
	return ffmpeg_go.Output(streams, output).Run()
}
