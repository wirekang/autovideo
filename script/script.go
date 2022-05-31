package script

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/wirekang/autovideo/inputer"
)

type Line struct {
	Text     string `json:"t"`
	Duration int    `json:"d"`
}

func ExtractThumbnail(filepath string) (string, []string, error) {
	v, err := os.ReadFile(filepath)
	if err != nil {
		return "", nil, fmt.Errorf("can't extract thumbnail: %w", err)
	}

	lines := strings.Split(string(v), "\n")
	if len(lines) < 3 {
		return "", nil, fmt.Errorf("file is too short")
	}

	return lines[0], lines[1:], nil
}

func SaveLines(lines []Line, filepath string) error {
	v, err := json.Marshal(lines)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, v, 0666)
}

func LoadLines(filepath string) ([]Line, error) {
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
		if line.Duration < 10 {
			return nil, fmt.Errorf("line %d millis is too short: %d", i, line.Duration)
		}
	}

	return lines, nil
}

func GenerateLines(lineStrings []string) ([]Line, error) {
	durations := make([]int, len(lineStrings))

	currentIndex := 0
	var startMillis, currentMillis int64
	startMillis = time.Now().UnixMilli()

	const width = 60

	borderString := color.BlueString(strings.Repeat("-", width))

	printInfo := func() {
		screen.Clear()
		var pre1, pre2, next1, next2 string

		if currentIndex > 1 {
			pre2 = color.HiBlackString(lineStrings[currentIndex-2])
			pre1 = color.HiBlackString(lineStrings[currentIndex-1])
		} else if currentIndex == 1 {
			pre2 = borderString
			pre1 = color.HiBlackString(lineStrings[currentIndex-1])
		} else if currentIndex == 0 {
			pre2 = ""
			pre1 = borderString
		}

		if currentIndex < len(lineStrings)-2 {
			next1 = color.HiBlackString(lineStrings[currentIndex+1])
			next2 = color.HiBlackString(lineStrings[currentIndex+2])
		} else if currentIndex == len(lineStrings)-2 {
			next1 = color.HiBlackString(lineStrings[currentIndex+1])
			next2 = borderString
		} else if currentIndex == len(lineStrings)-1 {
			next1 = borderString
			next2 = ""
		}

		text := lineStrings[currentIndex]

		fmt.Println(pre2)
		fmt.Println(pre1)
		fmt.Println(">> ", text)
		fmt.Println(next1)
		fmt.Println(next2)
		fmt.Println(currentMillis)
	}

	step := func() {
		currentMillis = time.Now().UnixMilli() - startMillis
		startMillis = time.Now().UnixMilli()
		durations[currentIndex] = int(currentMillis)
	}

	var stop func()

	onPerv := func() {
		step()
		currentIndex--
		if currentIndex == -1 {
			currentIndex = 0
		}
		printInfo()
	}

	onNext := func() {
		step()
		currentIndex++
		if currentIndex == len(lineStrings) {
			stop()
			return
		}
		printInfo()
	}

	stopChan := make(chan error)
	lines := make([]Line, len(lineStrings))

	onStop := func() {
		go func() {
			stopChan <- fmt.Errorf("input stopped")
		}()
	}

	inp := newInputer(onPerv, onNext, onStop)

	stop = func() {
		inp.Stop()
		screen.Clear()
		for i := range lines {
			lines[i] = Line{
				Text:     lineStrings[i],
				Duration: durations[i],
			}
		}
		stopChan <- nil
	}

	step()
	printInfo()
	go func() {
		_ = inp.Start()
	}()

	err := <-stopChan
	return lines, err
}

func newInputer(onPrev, onNext, onStop func()) *inputer.Inputer {
	return inputer.New(inputer.Option{
		OnPrev:  onPrev,
		OnNext:  onNext,
		OnStop:  onStop,
		PrevKey: 'q',
		NextKey: 'w',
		StopKey: 'r',
	})
}
