package main

import (
	"fmt"
	"sync"
)

//正确示例
func wg() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("wg1 finished")
}

//Tip:错误示例1, 不要把waitGroup传到go func() 里面,waitGroup会进行copy..里面跟外面不是同一个wg,会造成逻辑混乱
func errWg1() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg sync.WaitGroup) {
			fmt.Println(i)
			defer wg.Done()
		}(i, wg)
	}
	wg.Wait()
	fmt.Println("errWg1 finished")

}

//Tip: 错误示例2 add 不要放到go func()里面 waitGroup 可能提前结束wait
func errWg2() {
	var wg sync.WaitGroup

	for i := 0; i <= 1000; i++ {
		go func(i int) {
			wg.Add(1)
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("errWg2 finished")

}

//Tip: 错误示例3 用几个go func() 就add 几个。 waitGroup add 很大会有什么影响
func errWg3() {
	var wg sync.WaitGroup

	for i := 0; i <= 1000; i++ {
		wg.Add(100)
		go func(i int) {
			fmt.Println(i)
			defer wg.Done() // wg.add会减一 .所以不要太大，也不能负数
		}(i)
	}
	wg.Wait()
	fmt.Println("errWg2 finished")

}

func main() {
	//wg()
	//errWg1()
	//errWg2()
	errWg3()
}
