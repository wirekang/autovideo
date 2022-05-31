package inputer

import (
	"sync"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type Inputer struct {
	onPrev func()
	onNext func()
	onStop func()
	prevKey,
	nextKey,
	stopKey rune
	mutex  sync.Mutex
	isStop bool
}

type Option struct {
	OnPrev func()
	OnNext func()
	OnStop func()
	PrevKey,
	NextKey,
	StopKey rune
}

func New(o Option) *Inputer {
	return &Inputer{
		onPrev:  o.OnPrev,
		onNext:  o.OnNext,
		onStop:  o.OnStop,
		prevKey: o.PrevKey,
		nextKey: o.NextKey,
		stopKey: o.StopKey,
	}
}

func (i *Inputer) Stop() {
	i.mutex.Lock()
	i.isStop = true
	i.mutex.Unlock()
}

func (i *Inputer) Start() error {
	i.mutex.Lock()
	i.isStop = false
	i.mutex.Unlock()

	for {
		i.mutex.Lock()
		if i.isStop {
			i.mutex.Unlock()
			return nil
		}
		i.mutex.Unlock()
		color.Blue("Prev: %q  Next: %q               Stop: %q\n", i.prevKey, i.nextKey, i.stopKey)
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			return err
		}

		switch char {
		case i.stopKey:
			i.onStop()
			return nil
		case i.prevKey:
			i.onPrev()
		case i.nextKey:
			i.onNext()
		}
	}
}
