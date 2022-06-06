package script

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

func Parse(filePath string) (string, []string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	lines := strings.Split(string(data), "\n")
	lines = lo.Map(lines, func(line string, _ int) string { return strings.TrimSpace(line) })
	switch l := len(lines); l {
	case 0:
		err = fmt.Errorf("empty script file")
	case 1:
		fallthrough
	case 2:
		err = fmt.Errorf("no texts in script file")
	}
	if err != nil {
		return "", nil, err
	}

	title := lines[0]
	lines = lines[2:]
	var rst []string
	for _, line := range lines {
		if line != "" {
			rst = append(rst, line)
		}
	}
	return title, rst, nil
}
