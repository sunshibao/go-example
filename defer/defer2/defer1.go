package main

import "fmt"

func main() {
	fmt.Println("reciprocal")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
}

//即"先进后出"特性
