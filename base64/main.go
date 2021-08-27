package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	sEnc := "FE+xd/3zRL8ZPbycutdF5p45s3opbySZhCQu3T3qcUL26//uNHjX20zNQZmxd/RHina7SoZkOEv6L3+P1PpKPTweyYAkhQr4fg83s80M+N6rzS55cK8rw5zOnS1edkZhWwtV9mCY2fYBOWBeY2hQ2R+yVYUp0m02bpuzIKvPmRVFp+hVWgrqdQPvn3S+RXI1f0LLXt5VxC4efkP2bbpGSg=="
	fmt.Printf("enc=[%s]\n", sEnc)

	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		fmt.Printf("base64 decode failure, error=[%v]\n", err)
	} else {
		fmt.Printf("dec=[%s]\n", sDec)
	}
}
