package script

import (
	"encoding/json"
	"fmt"
	"os"
)

type Line struct {
	Text   string
	Millis int
}

func Parse(filepath string) ([]Line, error) {
	v, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("can't read script file: %w", err)
	}

	var lines []Line
	err = json.Unmarshal(v, &lines)
	if err != nil {
		return nil, err
	}

	if len(lines) < 3 {
		return nil, fmt.Errorf("len(lines) is too short: %d", len(lines))
	}

	for i, line := range lines {
		if line.Millis < 100 {
			return nil, fmt.Errorf("line %d millis is too short: %d", i, line.Millis)
		}
	}

	return lines, nil
}
