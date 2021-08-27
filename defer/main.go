package main

import "fmt"

func testa() {
	fmt.Println("aaaaaaaaaaaaaaaaa")
}

func testb(x int) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		//recover() //可以打印panic的错误信息
		//fmt.Println(recover())
		if err := recover(); err != nil { //产生了panic异常
			fmt.Println(err)
		}

	}() //别忘了(), 调用此匿名函数
	fmt.Println("bbbbbbbbbb11111111")

	var a [10]int
	a[x] = 111 //当x为20时候，导致数组越界，产生一个panic，导致程序崩溃

	fmt.Println("bbbbbbbbbb222222222")   //recover时，panic之前的还继续执行，之后的不会执行
}

func testc() {
	fmt.Println("cccccccccccccccccc")
}

func main() {
	testa()
	testb(20) //当值是1的时候，就不会越界，值是20的时候，就会越界报错。
	testc()
}