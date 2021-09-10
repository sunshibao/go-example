package main

import "fmt"

func main() {
	s := make([]int, 3, 3)
	f(s)
	fmt.Printf("%p\n", s)
	fmt.Println(len(s))
	fmt.Println(cap(s))
}
func f(s []int) {
	for range s {
		s = append(s, 2)
	}
	fmt.Printf("%p\n", s)
	fmt.Println(len(s))
	fmt.Println(cap(s))
}
