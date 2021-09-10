package main

import (
	"fmt"
)

func change(s ...string) {
	s[0] = "Go"
	s = append(s, "playground")
	fmt.Println(s)
}

func main() {
	//4 101
	//1 001

	fmt.Println(8 << 1)
	fmt.Println(8 >> 1)
	fmt.Println(8 | 1) //或
	fmt.Println(8 & 1) //与
	fmt.Println(uintptr(0))
	fmt.Println(uint64(0))
	fmt.Println(^uintptr(0))
	fmt.Println(^uint64(0))
	fmt.Println(^uintptr(0) >> 63)
	fmt.Println(^uint64(0) >> 63)
	fmt.Println(4 << (^uintptr(0) >> 63))
	fmt.Println(4 << (^uint64(0) >> 63))

	fmt.Println(^0)
	fmt.Println(^1)
	fmt.Println(^2)
	fmt.Println(^3)
	fmt.Println(1 ^ -3)

	// ^作一元运算符表示是按位取反，包括符号位在内
	// uint 无符号类型
	// uint 8 ----- 20 ------ ^0001 0100 = 1110 1011
	// uint 16 ----- 20 ------ ^0000 0000 0001 0100 = 1111 1111 1110 1011
	// uintptr = uint64 ----- 20 -----
	// ^0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0001 0100
	// = 1111 1111 1111 1111 1111 1111 1111 1111 1111 1111 1111 1111 1111 1111 1110 1011

	// int8 有符号类型 ---- 20 -----   一个有符号位的^操作为 这个数+1的相反数
	var a uint8 = 20
	var b uint16 = 20
	var c uintptr = 20
	var d int8 = 20
	fmt.Println(^a, ^b, ^c, ^d)
}
