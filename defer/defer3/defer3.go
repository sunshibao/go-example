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

//正确的答案是num1为1,num2为2，这两个变量并不受num1++、num2++的影响，因为defer将语句放入到栈中时，也会将相关的值拷贝同时入栈。
