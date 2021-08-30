package main

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesCbc(t *testing.T) {
	orig := "http://c.biancheng.net/golang/"
	key := "astaxie12798akljzmknm.ahkjkljl;k"
	fmt.Println("原文：", orig)
	encryptCode := AesEncryptCFB([]byte(orig), []byte(key))
	fmt.Println("密文：", base64.StdEncoding.EncodeToString(encryptCode)) //密文要用base64转换

	decryptCode := AesDecryptCFB(encryptCode, []byte(key))
	fmt.Println("解密结果：", string(decryptCode))
}
