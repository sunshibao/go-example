package main

import (
	"fmt"
)

func main() {
	a := 1
	b := 2
	defer calc(a, calc(a, b)) //1,4
	a = 0
	defer calc(a, calc(a, b)) //2,3
}

func calc(x, y int) int {
	fmt.Println(x, y, x+y)
	return x + y
}

// 参数先执行

// 123
// 022
// 022
// 134
