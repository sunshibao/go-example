package main

import "fmt"

func main() {
	p1 := new(int)
	*p1 = 1
	fmt.Println("p1", p1)
	fmt.Println("*p1", *p1)
	fmt.Println("&p1", &p1)

	s1 := new([]int)
	s2 := []*int{}
	s3 := []int{1, 2}

	s2 = append(s2, p1)
	fmt.Println("s1: ", s1)
	fmt.Printf("s2:%v\n", s2)
	fmt.Printf("s2:%T\n", s2)
	fmt.Println("s3: ", s3)

}
