package utils

import (
	"math/rand"
	"time"
)

// 包含上下限 [min, max]
func GetRandomWithAll(min, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max-min+1) + min)
}

// 不包含上限 [min, max)
func GetRandomWithMin(min, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max-min) + min)
}

// 不包含下限 (min, max]
func GetRandomWithMax(min, max int) int64 {
	var res int64
	rand.Seed(time.Now().UnixNano())
Restart:
	res = int64(rand.Intn(max-min+1) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}

// 都不包含 (min, max)
func GetRandomWithNo(min, max int) int64 {
	var res int64
	rand.Seed(time.Now().UnixNano())
Restart:
	res = int64(rand.Intn(max-min) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}
