package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func main() {

	key := []byte("YOGXKYgwxpnSgo6wujBnAqwKMS9DDlwJ")
	plaintext := "kp99ZsbvvTP9lFoaIs4VIwJy7dO1pEhIFkOFgpz5"

	saltByte, _ := hex.DecodeString("52fdfc072182654f163f5f0f")
	fmt.Println(saltByte)
	//iv := []byte("000000000000") //make([]byte, 12)
	encrypted, _ := aesGcmEncrypt(plaintext, key, saltByte)
	// kj+cPeGVoNvI6v7iLEBbzKXotDT9AQV+XYcicvGMKSn7/yLfwZCouS4e1Cji3I9ZvMYvrPRSVxZq6sdV0IzJH5D2XloeWRajgPZ6uzOCgGbDdxEf7ywf63KYss2KQK6JrFhkG4zlbSY=
	encryptedHex := hex.EncodeToString(encrypted)
	fmt.Println("encrypted", encryptedHex)
	encrypted, _ = hex.DecodeString(encryptedHex)
	decrypted, _ := aesGcmDecrypt(encrypted, key, saltByte)
	// This is a plain text which need to be encrypted by Java AES 256 GCM Encryption Algorithm
	fmt.Println("decrypted", string(decrypted))
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
