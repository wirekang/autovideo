package imageConcater

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/samber/lo"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wirekang/autovideo/fileutil"
	"github.com/wirekang/autovideo/line"
)

type ImageConcater struct {
	inputDir   string
	outputFile string
	filePrefix string
	lines      []line.Line
}

type Option struct {
	InputDir        string
	OutputFile      string
	ImageFilePrefix string
	Lines           []line.Line
}

func New(o Option) *ImageConcater {
	c := lo.Must(fileutil.Count(o.InputDir, false))
	if c == 0 {
		panic(fmt.Errorf("no files in %s", o.InputDir))
	}

	if c != len(o.Lines) {
		panic(fmt.Errorf("lines length <-> images count mismatch: %d <-> %d", len(o.Lines), c))
	}

	lo.Must0(os.MkdirAll(path.Dir(o.OutputFile), fs.ModeDir))

	return &ImageConcater{
		inputDir:   o.InputDir,
		outputFile: o.OutputFile,
		filePrefix: o.ImageFilePrefix,
		lines:      o.Lines,
	}
}

func (i *ImageConcater) ConcatImages() error {
	const name = "test.ffconcat"
	err := i.makeFFConcat(name)
	if err != nil {
		return err
	}

	return ffmpeg_go.Input(name).Output(i.outputFile).Run()
}

func (i *ImageConcater) makeFFConcat(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	var v []string

	v = append(v, "ffconcat version 1.0")
	for j, l := range i.lines {
		imagePath := fileutil.ImageName(i.filePrefix, j)
		imagePath = path.Join(i.inputDir, imagePath)
		v = append(v, fmt.Sprintf("file '%s'\nduration %.2f", imagePath, float64(l.Millis)/1000))
		_ = l
	}

	for _, s := range v {
		_, err = f.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}
	return f.Close()
}
