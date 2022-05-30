package fileutil

import (
	"os"
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
