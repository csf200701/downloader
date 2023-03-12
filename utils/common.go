package utils

import (
	"math"
	"strconv"
)

func CalcSize(size int64) string {
	var show float64
	var unit string
	if size < 1024 {
		show = float64(size)
		unit = "B"
	} else if size < int64(math.Pow(1024, 2)) {
		show = float64(size) / math.Pow(1024, 1)
		unit = "KB"
	} else if size < int64(math.Pow(1024, 3)) {
		show = float64(size) / math.Pow(1024, 2)
		unit = "MB"
	} else if size < int64(math.Pow(1024, 4)) { //千兆
		show = float64(size) / math.Pow(1024, 3)
		unit = "GB"
	} else if size < int64(math.Pow(1024, 5)) { //太字节
		show = float64(size) / math.Pow(1024, 4)
		unit = "TB"
	} else if size < int64(math.Pow(1024, 6)) { //拍字节
		show = float64(size) / math.Pow(1024, 5)
		unit = "PB"
	} else if size < int64(math.Pow(1024, 7)) { //艾字节
		show = float64(size) / math.Pow(1024, 6)
		unit = "EB"
	} else if size < int64(math.Pow(1024, 8)) { //泽字节
		show = float64(size) / math.Pow(1024, 7)
		unit = "ZB"
	} else if size < int64(math.Pow(1024, 9)) { //尧字节
		show = float64(size) / math.Pow(1024, 8)
		unit = "YB"
	} else if size < int64(math.Pow(1024, 10)) { //
		show = float64(size) / math.Pow(1024, 9)
		unit = "BB"
	}
	return strconv.FormatFloat(show, 'f', 2, 64) + unit
}
