package main

import (
	"fmt"
	"sync"
)

func wg()  {
	var wg sync.WaitGroup

	for i:=0;i<10;i++{
		wg.Add(1)
		go func(i int) {
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("wg1 finished")
}
// Tip:waitGroup 不要进行copy..里面跟外面不是同一个wg,会造成逻辑混乱
func errWg1()  {
	var wg sync.WaitGroup

	for i:=0;i<10;i++{
		wg.Add(1)
		go func(i int,wg sync.WaitGroup) {
			fmt.Println(i)
			defer wg.Done()
		}(i,wg)
	}
	wg.Wait()
	fmt.Println("errWg1 finished")

}

// Tip:waitGroup 可能提前结束wait
func errWg2()  {
	var wg sync.WaitGroup

	for i:=0;i<=1000;i++{
		go func(i int) {
			wg.Add(1)
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("errWg2 finished")

}

// Tip:waitGroup add 很大会有什么影响
func errWg3()  {
	var wg sync.WaitGroup

	for i:=0;i<=1000;i++{
		wg.Add(100)
		go func(i int) {
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("errWg2 finished")

}

func main()  {
	//wg()
	//errWg1()
	//errWg2()
	errWg3()
}