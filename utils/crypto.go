package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	fnv "hash/fnv"
	"math"

	"golang.org/x/crypto/md4"
)

func Hash(text string, hashInstance hash.Hash, isHex bool) string {
	if isHex {
		arr, _ := hex.DecodeString(text) // 十六进制字符串转为十六进制字节数组
		hashInstance.Write(arr)          // 写入哈希实例对象
	} else {
		hashInstance.Write([]byte(text)) // 将字符串转换为字节数组，写入哈希对象
	}
	bytes := hashInstance.Sum(nil)  // 哈希值追加到参数后面，只获取原始值，不用追加，用nil，返回哈希值字节数组
	return fmt.Sprintf("%x", bytes) // 格式化输出哈希值
}

func Md4(text string) string {
	md4Crypto := md4.New()
	return Hash(text, md4Crypto, false)
}

func Md5(text string) string {
	md5Crypto := md5.New()
	return Hash(text, md5Crypto, false)
}

func Sha256(text string) string {
	sha256Crypto := sha256.New()
	return Hash(text, sha256Crypto, false)
}

func Sha512(text string) string {
	sha512Crypto := sha512.New()
	return Hash(text, sha512Crypto, false)
}

func Short(text string) string {
	codes := []byte{'7', 'b', 'z', '2', 'e', '5', 'a', 'h', '3', 'm'}
	h := fnv.New32a()
	h.Write([]byte(text))
	n := uint64(h.Sum32())
	s := ""
	for i := 9; i >= 0; i-- {
		w := n / uint64(math.Pow(10, float64(i))) % 10
		s += string(codes[w])
	}

	return s
}
