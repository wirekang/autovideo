package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/samber/lo"
	"github.com/wirekang/autovideo/exe"
	"github.com/wirekang/autovideo/old/fileutil"
)

var DoNotDeleteTempFiles bool

var deletes struct {
	sync.Mutex
	filePaths []string
}

func ListenDeleteOnExit() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		Delete()
		os.Exit(0)
	}()
}

func Delete() {
	if DoNotDeleteTempFiles {
		return
	}

	fmt.Println("Delete temp files")
	deletes.Lock()
	for _, filePath := range deletes.filePaths {
		tryRemove(filePath)
	}
	deletes.filePaths = []string{}
	deletes.Unlock()
}

func DeleteOnExit(filePath string) {
	deletes.Lock()
	deletes.filePaths = append(deletes.filePaths, filePath)
	deletes.Unlock()
}

func ForEachLines(filePath string, cb func(index int, line string) error) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	s := bufio.NewScanner(f)
	var i int
	for s.Scan() {
		err = cb(i, s.Text())
		if err != nil {
			return err
		}

		i++
	}
	return s.Err()
}

func MkdirAllReset(dirPath string) error {
	_ = os.RemoveAll(dirPath)
	return os.MkdirAll(dirPath, os.ModeDir)
}

func modTime(entry fs.DirEntry) time.Time {
	i, _ := entry.Info()
	return i.ModTime()
}

func ListByModTime(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return modTime(entries[i]).Before(fileutil.ModTime(entries[j]))
	})

	paths := lo.Map(entries, func(e os.DirEntry, _ int) string {
		return path.Join(dirPath, e.Name())
	})

	return paths, nil
}

func tryRemove(filepath string) {
	_ = os.RemoveAll(filepath)
	_ = os.Remove(filepath)
}

func Index(dirPath, prefix string, index int, extension string) string {
	return path.Join(dirPath, fmt.Sprintf("%s%03d.%s", prefix, index, extension))
}

func MediaLengthSeconds(filePath string) (float64, error) {
	n, err := exe.Run(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		filePath,
	)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(n, 32)
}

func MediaLengthSecondsAll(filePaths []string) ([]float64, error) {
	rst := make([]float64, len(filePaths))
	for i, filePath := range filePaths {
		l, err := MediaLengthSeconds(filePath)
		if err != nil {
			return nil, err
		}
		rst[i] = l
	}

	return rst, nil
}

func WriteFFConcat(filePaths []string, durations []float64, outputFilePath string) error {
	if durations != nil && len(filePaths) != len(durations) {
		return fmt.Errorf("length mismatch: filePaths(%d) <> durations(%d)", len(filePaths), len(durations))
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	var lines []string
	lines = append(lines, "ffconcat version 1.0")
	for i, filePath := range filePaths {
		lines = append(lines, fmt.Sprintf("file '%s'", filePath))
		if durations != nil {
			lines = append(lines, fmt.Sprintf("duration %.2f", durations[i]))
		}
	}
	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	DeleteOnExit(outputFilePath)
	return f.Close()
}
