package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("参数错误,rootKey.go 后面只能跟一个string类型的参数")
		os.Exit(1)
	}
	if "ppp" == os.Args[1] {
		fmt.Println("pppp")
	} else {
		fmt.Println("aaaaa")
	}
}
