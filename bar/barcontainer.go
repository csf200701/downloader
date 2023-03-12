package bar

import (
	"fmt"
	"os"
	"sync"

	term "github.com/nsf/termbox-go"
)

type BarContainer struct {
	sync.Mutex
	bars     []*Bar
	title    string
	nextLine int
	maxLine  int
	width    int
	minLine  int
}

func New(title string) *BarContainer {
	if !term.IsInit {
		term.Init()
	}
	// _, minLine, _ := termWidth()
	minLine := 0
	term.SetCursor(0, minLine)
	term.SetOutputMode(term.Output256)
	fmt.Printf("\r" + title)
	b := &BarContainer{nextLine: 0, title: title, minLine: minLine}
	b.setWidth()
	go func() {
		for {
			e := term.PollEvent()
			if e.Key == term.KeyCtrlC || e.Type == term.EventInterrupt || e.Type == term.EventError {
				b.Lock()
				defer b.Unlock()
				//term.Close()
				term.Flush()
				os.Exit(0)
				term.SetCursor(0, minLine)
			} else if e.Type == term.EventResize {
				b.setWidth()
				term.SetCursor(0, minLine)
				fmt.Printf(title)
			}
		}
	}()

	return b
}

func (bc *BarContainer) NewBar(prepend string, total int64) *Bar {
	if total <= 0 {
		return nil
	}
	bc.nextLine++
	b := newBar(bc, prepend, total, bc.nextLine)
	bc.bars = append(bc.bars, b)
	return b
}

var termWidth = func() (width int, height int, err error) {
	w, h := term.Size()
	return w, h, nil
}

func (bc *BarContainer) setWidth() {
	width, _, err := termWidth()
	if err != nil {
		width = 50
	} else {
		width = width - 80
	}
	bc.width = width
}
