package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {

	key := []byte("7UkX6l@bp8r@#y0Zqwe7lkS1cqL*1%Tn")
	plaintext := "{\"head\":\"{\\\"mcd\\\":101151,\\\"areaId\\\":\\\"RU\\\",\\\"acToken\\\":\\\"CgB6e3x9JN8mdJ\\\\\\/PbQXXWy8hXlAXHimj2e+PJf9pjxPx2e7T+H95R+y2gbx0krAsXLWEXR3PMSRF7LFFW1Z64yj\\\\\\/lxJaQUBuHBNnVBuiVS5t7ms=\\\"}\",\"body\":\"{\\\"tInfo\\\":{\\\"androidId\\\":\\\"185a796aba13bc4b\\\",\\\"apkVerName\\\":\\\"16.0.1.109\\\",\\\"apkVer\\\":160001109,\\\"chId\\\":\\\"HONOR_01\\\",\\\"hman\\\":\\\"HONOR\\\",\\\"htype\\\":\\\"NTN-LX1\\\",\\\"isDark\\\":false,\\\"language\\\":\\\"zh_cn\\\",\\\"netType\\\":3,\\\"openId\\\":\\\"MDFAMjExMDAwMjIxQDhkMWVhMWZlMDc1NGVmODNkYmU4ZGE0ZTVlZjY1YTcwQDQwODE3OGRmYTRkZjUxNGViaMWViaZDRmNTA4ZmQxNmVkYzNhMGZlNzRkNzU0ZTczYmY1OGIxY2JhY2U\\\",\\\"osVer\\\":\\\"11\\\",\\\"pName\\\":\\\"com.hihonor.gamecenter\\\",\\\"randomId\\\":\\\"273539791636628582073\\\"}}\"}"
	iv := []byte("560472785304") //make([]byte, 12)
	encrypted, _ := AESGCMEncrypter(plaintext, key, iv)
	// kj+cPeGVoNvI6v7iLEBbzKXotDT9AQV+XYcicvGMKSn7/yLfwZCouS4e1Cji3I9ZvMYvrPRSVxZq6sdV0IzJH5D2XloeWRajgPZ6uzOCgGbDdxEf7ywf63KYss2KQK6JrFhkG4zlbSY=
	fmt.Println("encrypted", encrypted)
	//encrypted="LSpsHZiuTLKZEQjLUC9I+UwiZKvkUhQ0VmZRCvLhCXSd/MsGchK0r3KYTj7Z5MW/i4lyN5vz57S7CebP/+IHbxo/D6psZWOIeQEWgFIr4YMWZK8J9fbD09zHG3xMaw8gPHQtHqlLUpPVc2RM44lPQaCbEKHFNkfgJ3zS1rV9AqweG5xSW1L/xHWfLpKZvZ3htwoaj8sZ6kdnYk36MHNRuCkzE1xhPWEzLCiFxu6E5A/O+aM+QEnEpuZbhrO9qx/hPx0chnDTe9fccnkAjrF1mL+aykneZCZdl4thtWmff+fjVCk4tnIcDxMq5dm5eUYfs14oDIIA3KEGVlKnjFOZRGo6SruQ8Q99mhBvaBPLOiwAzD+g8i8Nwd1WRtUotz8B/ALj3w/A3WOakMeRnwnzJtrYgnGrWtA63/4FdAX4MZs3Mt3kRocDvFylvLYB9LJtV9nl90AUoxD2yBXXbekywWXrIyrUNl0C4A=="
	//encrypted = "LSpsHZiuTLKZEQjLUC9I+UwiZKvkUhQ0VmZRCvLhCXSd/MsGchK0r3KYTj7Z5MW/i4lyN5vz57S7CebP/+IHbxo/D6psZWOIeQEWgFIr4YMWZK8J9fbD09zHG3xMaw8gPHQtHqlLUpPVc2RM44lPQaCbEKHFNkfgJ3zS1rV9AqweG5xSW1L/xHWfLpKZvZ3htwoaj8sZ6kdnYk36MHNRuCkzE1xhPWEzLCiFxu6E5A/O+aM+QEnEpuZbhrO9qx/hPx0chnDTe9fccnkAjrF1mL+aykneZCZdl4thtWmff+fjVCk4tnIcDxMq5dm5eUYfs14oDIIA3KEGVlKnjFOZRGo6SruQ8Q99mhBvaBPLOiwAzD+g8i8Nwd1WRtUotz8B/ALj3w/A3WOakMeRnwnzJtrYgnGrWtA63/4FdAX4MZs3Mt3kRocDvFylvLYB9LJtV9nl90AUoxD2yBXXbekywWXrIyrUNl0C4A=="
	//encrypted="LSpsHZiuTLKZEQjLUC9I+UwiZ6vkUhQ0VmZRCvLhCXSd/MsGchK0r3KYH0DXp8SEgJs1Y+Wqgcr4JbT+tL1RUVdhLdV6GzqwOgE6shZ+4YNwPd93tv3V09zcPFpMUnE4WmZnde0bRtWIOUtK9sweQ/2kbrm1SDmje0um2qlHBrMNJtI0JUr+xhmTdb32+K+nnxQCrc8H2gA5DUOcIjxl0yQcYg5SaHVeWAG29a3QxQ+fNyi36zGwVpb5gh0CUAiiiA=="
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
