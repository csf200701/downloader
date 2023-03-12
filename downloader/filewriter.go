package downloader

import (
	"downloader/bar"
	"math"
	"os"
	"time"
)

type FileWriter struct {
	bc         *bar.BarContainer
	bar        *bar.Bar
	prepend    string
	start      int64
	end        int64
	size       int64
	offset     int64
	file       *os.File
	secTime    int
	ratePerSec int
	sizePerSec int
}

func NewFileWriter(bc *bar.BarContainer, bar *bar.Bar, prepend string, start, end, size int64) *FileWriter {
	file, _ := os.OpenFile(prepend, os.O_CREATE|os.O_RDWR, 0644)
	var left int64 = 0
	fileInfo, _ := file.Stat()
	left += fileInfo.Size()
	fw := &FileWriter{bar: bar, prepend: prepend, file: file, start: start, end: end, size: size, offset: left, secTime: int(time.Now().UnixNano())}
	fw.bar.Add(left, 0)
	return fw
}

func (w *FileWriter) Write(p []byte) (int, error) {
	n := len(p)
	w.file.WriteAt(p, w.offset)
	w.offset += int64(n)
	curTime := int(time.Now().UnixNano())
	w.sizePerSec += n
	if curTime-w.secTime >= int(math.Pow10(9)) {
		w.ratePerSec = int(math.Pow10(9)) * w.sizePerSec / (curTime - w.secTime)
		w.secTime = curTime
		w.sizePerSec = 0
	}
	w.bar.Add(int64(n), w.ratePerSec)
	return n, nil
}

func (w *FileWriter) Close() error {
	if w.file == nil {
		return nil
	}
	return w.file.Close()
}
