package main

import "fmt"

func main() {
	anonymous := Anonymous()
	name := HasName()
	fmt.Println(anonymous)
	fmt.Println(name)
}

// 匿名函数
func Anonymous() int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2 value is ", i)
	}()

	defer func() {
		i++
		fmt.Println("defer1 value is ", i)
	}()

	return i
}

//defer 函数能够在 return 语句对返回值赋值之后，继续对返回值进行操作

func HasName() (j int) {
	defer func() {
		j++
		fmt.Println("defer2 in value", j)
	}()

	defer func() {
		j++
		fmt.Println("defer1 in value", j)
	}()

	return j

}
