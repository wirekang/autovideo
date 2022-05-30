package imageConcater

import (
	"fmt"
	"os"
	"path"

	"github.com/samber/lo"
	"github.com/wirekang/autovideo/ffmpegutil"
	"github.com/wirekang/autovideo/fileutil"
	"github.com/wirekang/autovideo/script"
)

type ImageConcater struct {
	inputDir   string
	outputFile string
	filePrefix string
	lines      []script.Line
}

type Option struct {
	InputDir        string
	OutputFile      string
	ImageFilePrefix string
	Lines           []script.Line
}

func New(o Option) *ImageConcater {
	c := lo.Must(fileutil.Count(o.InputDir, false))
	if c == 0 {
		panic(fmt.Errorf("no files in %s", o.InputDir))
	}

	if c != len(o.Lines) {
		panic(fmt.Errorf("lines length <-> images count mismatch: %d <-> %d", len(o.Lines), c))
	}

	lo.Must0(fileutil.MkdirAll(path.Dir(o.OutputFile)))
	return &ImageConcater{
		inputDir:   o.InputDir,
		outputFile: o.OutputFile,
		filePrefix: o.ImageFilePrefix,
		lines:      o.Lines,
	}
}

func (i *ImageConcater) Concat() error {
	fileutil.TryRemove(i.outputFile)
	const name = "images.ffconcat"
	defer func() {
		fileutil.TryRemove(name)
	}()

	err := i.makeFFConcat(name)
	if err != nil {
		return err
	}

	return ffmpegutil.InputOutput(name, i.outputFile)
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
	}

	for _, s := range v {
		_, err = f.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}
	return f.Close()
}
