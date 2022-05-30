package fileutil

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

func Count(path string, includeDir bool) (int, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	count := 0

	for _, entry := range entries {
		if includeDir || !entry.IsDir() {
			count++
		}
	}

	return count, nil
}

func ImageName(prefix string, index int) string {
	return fmt.Sprintf("%s%04d.png", prefix, index)
}

func MkdirAll(dir string) error {
	_ = os.Remove(dir)
	return os.MkdirAll(dir, os.ModeDir)
}

func ModTime(entry fs.DirEntry) time.Time {
	i, _ := entry.Info()
	return i.ModTime()
}
