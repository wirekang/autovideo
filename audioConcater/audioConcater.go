package audioConcater

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/samber/lo"
	"github.com/wirekang/autovideo/ffmpegutil"
	"github.com/wirekang/autovideo/fileutil"
)

type AudioConcater struct {
	inputDir   string
	outputFile string
}

type Option struct {
	InputDir   string
	OutputFile string
}

func New(o Option) *AudioConcater {
	c := lo.Must(fileutil.Count(o.InputDir, false))
	if c == 0 {
		panic(fmt.Errorf("no files in %s", o.InputDir))
	}

	lo.Must0(fileutil.MkdirAll(path.Dir(o.OutputFile)))
	return &AudioConcater{
		inputDir:   o.InputDir,
		outputFile: o.OutputFile,
	}
}

func (i *AudioConcater) Concat() error {
	fileutil.TryRemove(i.outputFile)
	const name = "audios.ffconcat"
	defer func() {
		fileutil.TryRemove(name)
	}()

	err := i.makeFFConcat(name)
	if err != nil {
		return err
	}

	return ffmpegutil.InputOutput(name, i.outputFile)
}

func (i *AudioConcater) makeFFConcat(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(i.inputDir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return fileutil.ModTime(files[i]).Before(fileutil.ModTime(files[j]))
	})

	var v []string

	v = append(v, "ffconcat version 1.0")
	for _, file := range files {
		audioPath := path.Join(i.inputDir, file.Name())
		v = append(v, fmt.Sprintf("file '%s'", audioPath))
	}

	for _, s := range v {
		_, err = f.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}
	return f.Close()
}
