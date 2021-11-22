package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sunshibao/go-utils/util/gconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("参数错误,只能跟一个参数")
		os.Exit(1)
	}
	// 多少位
	rootLen := os.Args[1]

	saltByte := []byte(getPlainWork(gconv.Int(rootLen)))

	encryptedHex := hex.EncodeToString(saltByte)

	fmt.Println("rootKey:" + strings.ToUpper(encryptedHex))
}

//随机生成plainWork
func getPlainWork(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var src = rand.NewSource(time.Now().UnixNano())

	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
