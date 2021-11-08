package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	key := []byte("7UkX6l@bp8r@#y0Zqwe7lkS1cqL*1%Tn")
	plaintext := "kp99ZsbvvTP9lFoaIs4VIwJy7dO1pEhIFkOFgpz5"
	iv := []byte("000000000000") //make([]byte, 12)
	encrypted, _ := AESGCMEncrypter(plaintext, key, iv)
	// kj+cPeGVoNvI6v7iLEBbzKXotDT9AQV+XYcicvGMKSn7/yLfwZCouS4e1Cji3I9ZvMYvrPRSVxZq6sdV0IzJH5D2XloeWRajgPZ6uzOCgGbDdxEf7ywf63KYss2KQK6JrFhkG4zlbSY=
	fmt.Println("encrypted", encrypted)
	encrypted = "HnFMIsqgBNiODDytS3geht7sUWOA3Y0a8aau9FVVPOg="
	decrypted, _ := AESGCMDecrypter(encrypted, key, iv)
	// This is a plain text which need to be encrypted by Java AES 256 GCM Encryption Algorithm
	fmt.Println("decrypted", string(decrypted))
}

func AESGCMEncrypter(src string, key, iv []byte) (encrypted string, err error) {
	plaintext := []byte(src)
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

	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)
	encrypted = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

func AESGCMDecrypter(crypt string, key, iv []byte) (decrypted []byte, err error) {
	ciphertext, errBase64 := base64.StdEncoding.DecodeString(crypt)
	if errBase64 != nil {
		err = errBase64
		return
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

	decrypted, err = aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return
	}

	return decrypted, err
}
