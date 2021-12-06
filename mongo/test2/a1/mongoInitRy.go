package main

import "fmt"

func main() {
	fmt.Println(Sum(1, 2))
}

func Sum(num1, num2 int) int {
	defer fmt.Println("num1:", num1)
	defer fmt.Println("num2:", num2)
	num1++
	num2++
	return num1 + num2
}
