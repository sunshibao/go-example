package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("参数错误,只能跟两个参数")
		os.Exit(1)
	}
	// 要解析的明文
	plainText := os.Args[1]
	// 密文
	cipherText := os.Args[2]

	split := strings.Split(cipherText, ":")
	if len(split) != 2 {
		fmt.Println("第二个参数错误。格式[iv:密文]")
		os.Exit(1)
	}
	iv, err := hex.DecodeString(split[0])
	if err != nil {
		fmt.Println("iv 解析失败")
	}
	//秘钥 - 反解析密文
	key, err := hex.DecodeString(split[1])
	if err != nil {
		fmt.Println("秘钥 解析失败")
	}
	rootKey := getRootKey()
	workKey, err := aesGcmDecrypt(key, []byte(rootKey), iv)
	if err != nil {
		fmt.Println("workKey 生成失败")
	}

	saltByte := []byte(getPlainWork2(12))

	encrypted, _ := aesGcmEncrypt(plainText, workKey, saltByte)
	encryptedHex := hex.EncodeToString(encrypted)

	encrypted, _ = hex.DecodeString(encryptedHex)
	decrypted, _ := aesGcmDecrypt(encrypted, workKey, saltByte)
	fmt.Println("decrypted", string(decrypted))
	fmt.Println(hex.EncodeToString(saltByte) + ":" + encryptedHex)
}

//数据加密
func aesGcmEncrypt(src string, key []byte, iv []byte) ([]byte, error) {
	plaintext := []byte(src)
	block, errNewCipher := aes.NewCipher(key)
	if errNewCipher != nil {
		return nil, errNewCipher
	}

	aesgcm, errNewGCM := cipher.NewGCM(block)
	if errNewGCM != nil {
		return nil, errNewGCM
	}

	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)
	return ciphertext, nil
}

//数据解密
func aesGcmDecrypt(cryptByte []byte, key []byte, iv []byte) (decrypted []byte, err error) {
	if key == nil || iv == nil {
		return nil, errors.New("key fail")
	}

	block, errNewCipher := aes.NewCipher(key)
	if errNewCipher != nil {
		err = errNewCipher
		return
	}

	aesgcm, errNewGCM := cipher.NewGCM(block)
	if errNewGCM != nil {
		err = errNewGCM
		return
	}

	decrypted, err = aesgcm.Open(nil, iv, cryptByte, nil)
	if err != nil {
		return
	}

	return decrypted, err
}

//随机生成plainWork
func getPlainWork2(n int) string {
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

//读取文件
func ReadFile(fileName string) string {
	bytes, _ := ioutil.ReadFile(fileName)
	return string(bytes)
}

//string 转有符合的int
func strToInt(str []byte) []int {
	result := []int{}
	for _, v := range str {
		t := int(v)
		if v > 255/2 {
			t = int(v) - 256
		}
		result = append(result, t)
	}
	return result
}

func getRootKey() string {
	var thirdKey = "93093ADA767B6A59BD35E6208960119F4CFC1839544282183B5050BCAF83F9D995B3E214764F1970E45945B7FAEC6AA01BD7C965F2B611A04F325A04471711F9E78EFCC38831ACE419D1D3D76F04DA3A13314A58CF6B8F85CF346D67276B60E77B3AB980E6453F7C71DB4519083996CABC0713F0399C99D03F0E10D836305E2F"

	first := strings.Trim(ReadFile("../honor/honord/k.txt"), "")
	second := strings.Trim(ReadFile("../honor/honorw/k.txt"), "")
	salt := strings.Trim(ReadFile("../honor/honors/s.txt"), "")

	fmt.Println(first)
	c1, _ := hex.DecodeString(first)
	c2, _ := hex.DecodeString(second)
	c3, _ := hex.DecodeString(thirdKey)
	sByte, _ := hex.DecodeString(salt)

	c1Int := strToInt(c1)
	c2Int := strToInt(c2)
	c3Int := strToInt(c3)
	//sInt := strToInt(sByte)

	lenght := len(c1Int)
	if lenght > len(c2Int) {
		lenght = len(c2Int)
	}
	if lenght > len(c3Int) {
		lenght = len(c3Int)
	}

	combinedByte := []byte{}
	for i := 0; i < lenght; i++ {
		t := c1[i] ^ c2[i] ^ c3[i]
		combinedByte = append(combinedByte, uint8(t))
	}

	rootKeyByte := pbkdf2.Key(combinedByte, sByte, 10000, 16, sha256.New)
	rootKey := hex.EncodeToString(rootKeyByte)
	return rootKey
}
