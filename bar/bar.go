package bar

import (
	"downloader/utils"
	"fmt"
	"strconv"
	"sync"
	"time"

	term "github.com/nsf/termbox-go"
)

type Bar struct {
	mu         sync.Mutex
	bc         *BarContainer
	prepend    string
	rate       string    // 进度条
	current    int64     // 当前进度位置
	total      int64     // 总进度
	start      time.Time // 开始时间
	line       int       //第几行
	left       rune
	right      rune
	fill       rune
	head       rune
	empty      rune
	ratePerSec int
}

func newBar(bc *BarContainer, prepend string, total int64, line int) *Bar {
	b := &Bar{
		bc:      bc,
		prepend: prepend,
		current: 0,
		total:   total,
		start:   time.Now(),
		line:    line,
		left:    '[',
		right:   ']',
		fill:    '=',
		head:    '>',
		empty:   '-',
		// left:    '╢',
		// right:   '╟',
		// fill:    '▌',
		// head:    '▌',
		// empty:   '░',
	}
	b.load()
	return b
}

var ss int = 0

func (bar *Bar) load() {
	bar.bc.Lock()
	defer bar.bc.Unlock()
	bar.genRate()
	currentShow := utils.CalcSize(bar.current)
	totalShow := utils.CalcSize(bar.total)
	percent := int(bar.getPercent() * 100)
	term.SetCursor(0, bar.bc.minLine+bar.line)
	// term.SetFg(0, bar.bc.minLine+bar.line, term.ColorBlue)
	if bar.current == bar.total {
		fmt.Printf("\r"+bar.prepend+": "+string(bar.left)+"%-"+strconv.Itoa(bar.bc.width)+"s"+string(bar.right)+"% 4d%%  %-10s %-18s %-18s", bar.rate, percent, bar.getTime(), currentShow+"/"+totalShow, "完成")
	} else {
		fmt.Printf("\r"+bar.prepend+": "+string(bar.left)+"%-"+strconv.Itoa(bar.bc.width)+"s"+string(bar.right)+"% 4d%%  %-10s %-18s %-18s", bar.rate, percent, bar.getTime(), currentShow+"/"+totalShow, "进行中（"+utils.CalcSize(int64(bar.ratePerSec))+"/s）")
	}
	term.SetCursor(0, bar.bc.minLine+bar.bc.nextLine+1)
	// term.HideCursor()
	term.Flush()
}

func (bar *Bar) Reset(current int64) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	if current > bar.total {
		return
	}
	bar.current = current
	bar.load()

}

func (bar *Bar) Add(i int64, ratePerSec int) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	if bar.current+i > bar.total {
		return
	}
	bar.current += i
	bar.ratePerSec = ratePerSec
	bar.load()
}

func (bar *Bar) getPercent() float64 {
	return float64(bar.current) / float64(bar.total)
}

func (bar *Bar) getPos() int {
	return int(bar.getPercent() * float64(bar.bc.width))
}

func (bar *Bar) genRate() {
	curPos := bar.getPos()
	bar.rate = ""
	for i := 0; i < curPos; i++ {
		bar.rate += string(bar.fill)
	}
	plus := 0
	if bar.current > 0 && bar.current < bar.total {
		bar.rate += string(bar.head)
		plus++
	}
	for i := 0; i < bar.bc.width-curPos-plus; i++ {
		bar.rate += string(bar.empty)
	}
}

func (bar *Bar) getTime() (s string) {
	u := time.Now().Sub(bar.start).Seconds()
	h := int(u) / 3600
	m := int(u) % 3600 / 60
	if h > 0 {
		s += strconv.Itoa(h) + "h "
	}
	if h > 0 || m > 0 {
		s += strconv.Itoa(m) + "m "
	}
	s += strconv.Itoa(int(u)%60) + "s"
	return
}
