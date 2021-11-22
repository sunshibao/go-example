package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

/*
PlainrootKey: 332029f20ee529de9739c61012ba2112
PlainWorkKey: 0450E5A7A9A888A1C7C8F1034ADDF5708065A74F837F71DC75DC5F28460AE2B2
Encrypt result key: AC40C1268C22866FD73599CB:96F576A630DB8D422970133106BEF300A6EAE247BD519E6063E3A67238A9C778BF2C06A859C0BB157EE625411FE5A31D
Encrypt result mac: DC6B6EF8CA78042E1B3EF5CCAA624A24771338012A4E9C626651C1FA7E148130
*/

func main() {

	//获得根秘钥
	rootKey := getRootKey()

	fmt.Println(rootKey)

	//随机salt
	saltByte := getPlainWork(12)

	//随机生成plainWork
	plainWork := getPlainWork(32)

	saltHex := hex.EncodeToString([]byte(saltByte))

	//把plainWork加密保存
	plainWorkSecret, err := aesGCMEncrypt(plainWork, []byte(rootKey), []byte(saltByte))
	if err != nil {
		fmt.Println(err)
		return
	}
	plainWorkSecretHex := hex.EncodeToString(plainWorkSecret)
	//写入配置中保存
	fmt.Println("workKey:", saltHex+":"+plainWorkSecretHex)

}

//读取文件
func ReadFile(fileName string) string {
	bytes, _ := ioutil.ReadFile(fileName)
	return string(bytes)
}

func getRootKey() string {
	var thirdKey = "93093ADA767B6A59BD35E6208960119F4CFC1839544282183B5050BCAF83F9D995B3E214764F1970E45945B7FAEC6AA01BD7C965F2B611A04F325A04471711F9E78EFCC38831ACE419D1D3D76F04DA3A13314A58CF6B8F85CF346D67276B60E77B3AB980E6453F7C71DB4519083996CABC0713F0399C99D03F0E10D836305E2F"
	first := strings.Trim(ReadFile("aes/honor/honord/k.txt"), "")
	second := strings.Trim(ReadFile("aes/honor/honorw/k.txt"), "")
	salt := strings.Trim(ReadFile("aes/honor/honors/s.txt"), "")

	c1, _ := hex.DecodeString(first)
	c2, _ := hex.DecodeString(second)
	c3, _ := hex.DecodeString(thirdKey)
	sByte, _ := hex.DecodeString(salt)

	lenght := len(c1)
	if lenght > len(c2) {
		lenght = len(c2)
	}
	if lenght > len(c3) {
		lenght = len(c3)
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

// int8转btye
func intToByte(data []int32) []int32 {
	result := []rune{}
	for _, v := range data {
		if v < 0 {
			v = v + 65535
		}
		result = append(result, int32(v))
	}
	return result
}

//获取随机数 作为iv
func getRandIv() []byte {
	iv := make([]byte, 12)
	rand.Read(iv)
	return iv
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

// hmacsha256验证
func hmacSha256(src, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return hex.EncodeToString(m.Sum(nil))
}

// sha256验证
func sHA256Str(src string) string {
	h := sha256.New()
	h.Write([]byte(src)) // 需要加密的字符串为
	// fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil))) // 输出加密结果
	return hex.EncodeToString(h.Sum(nil))
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(ciphertext []byte) []byte {
	lenght := len(ciphertext)
	unpadding := int(len(ciphertext) - 1)
	return ciphertext[:(lenght - unpadding)]
}

//数据加密
func aesGCMEncrypt(src string, key []byte, salt []byte) ([]byte, error) {
	iv := salt // []byte("000000000000") //make([]byte, 12)
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
