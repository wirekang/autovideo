package script

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/wirekang/autovideo/inputer"
)

type Line struct {
	Text   string
	Millis int
}

// func Parse(filepath string) ([]Line, error) {
// 	v, err := os.ReadFile(filepath)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't read script file: %w", err)
// 	}
//
// 	var lines []Line
// 	err = json.Unmarshal(v, &lines)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(lines) < 3 {
// 		return nil, fmt.Errorf("len(lines) is too short: %d", len(lines))
// 	}
//
// 	for i, line := range lines {
// 		if line.Millis < 10 {
// 			return nil, fmt.Errorf("line %d millis is too short: %d", i, line.Millis)
// 		}
// 	}
//
// 	return lines, nil
// }

func Generate(txtFilePath string) ([]Line, error) {
	bytes, err := os.ReadFile(txtFilePath)
	if err != nil {
		return nil, err
	}

	texts := strings.Split(string(bytes), "\n")
	if len(texts) < 3 {
		return nil, fmt.Errorf("%s is too short", txtFilePath)
	}

	durations := make([]int, len(texts))

	currentIndex := 0
	var startMillis, currentMillis int64
	startMillis = time.Now().UnixMilli()

	const width = 60

	borderString := color.BlueString(strings.Repeat("-", width))

	printInfo := func() {
		screen.Clear()
		var pre1, pre2, next1, next2 string

		if currentIndex > 1 {
			pre2 = color.HiBlackString(texts[currentIndex-2])
			pre1 = color.HiBlackString(texts[currentIndex-1])
		} else if currentIndex == 1 {
			pre2 = borderString
			pre1 = color.HiBlackString(texts[currentIndex-1])
		} else if currentIndex == 0 {
			pre2 = ""
			pre1 = borderString
		}

		if currentIndex < len(texts)-2 {
			next1 = color.HiBlackString(texts[currentIndex+1])
			next2 = color.HiBlackString(texts[currentIndex+2])
		} else if currentIndex == len(texts)-2 {
			next1 = color.HiBlackString(texts[currentIndex+1])
			next2 = borderString
		} else if currentIndex == len(texts)-1 {
			next1 = borderString
			next2 = ""
		}

		text := texts[currentIndex]

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
		if currentIndex == len(texts) {
			stop()
			return
		}
		printInfo()
	}

	stopChan := make(chan error)
	lines := make([]Line, len(texts))

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
				Text:   texts[i],
				Millis: durations[i],
			}
		}
		stopChan <- nil
	}

	step()
	printInfo()
	go func() {
		_ = inp.Start()
	}()

	err = <-stopChan
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
