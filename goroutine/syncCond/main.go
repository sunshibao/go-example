package main

import (
	"fmt"
	"sync"
	"time"
)

//错误示例
func syncCondErr() {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	go func() {
		time.Sleep(1 * time.Second)
		c.Broadcast()
	}()
	m.Lock()
	time.Sleep(2 * time.Second)
	c.Wait()
}

//正确示例
func syncCondExplain() {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	// Tip: 主协程先获得锁
	c.L.Lock()
	go func() {
		// Tip:协程一开始无法获得锁
		c.L.Lock()
		defer c.L.Unlock()
		fmt.Println("3.该协程获得了锁")
		time.Sleep(2 * time.Second)
		// Tip: 通过notify进行广播通知
		c.Broadcast()
		fmt.Println("4.该协程执行完毕，即将执行defer中的解锁操作")
	}()

	fmt.Println("1.主协程获得锁")
	time.Sleep(1 * time.Second)
	fmt.Println("2.主协程依旧抢占着锁获得锁")
	//Tip:看一下Wait的大致实现，可以了解到，他是先释放锁，直到收到notify，又进行加锁
	c.Wait()
	//Tip:记得释放锁
	c.L.Unlock()

	fmt.Println("down")
}

func syncCond() {
	lock := sync.Mutex{}
	cond := sync.NewCond(&lock)
	for i := 0; i < 5; i++ {
		go func(i int) {
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait()
			fmt.Printf("No.{%d} Goroutine Receive \n", i)
		}(i)
	}
	time.Sleep(time.Second)
	cond.Broadcast() // 广播
	//cond.Signal()  单播
	time.Sleep(time.Second)
}

func main() {
	syncCond()
}
