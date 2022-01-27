package main

import "fmt"

func main() {
	fmt.Println(Test1())
	fmt.Println(Test2())
	fmt.Println(Test3())

}

func Test1() (r int) {
	i := 1
	defer func() {
		i = i + 1
	}()
	return i
}

func Test2() (r int) {
	defer func(r int) {
		r = r + 2
	}(r)
	return 2
}

func Test3() (r int) {
	defer func(r *int) {
		*r = *r + 2
	}(&r)
	return 2
}
