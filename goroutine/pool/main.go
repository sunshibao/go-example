package main

import (
	"fmt"
	"sync"
	"time"
)

type structR6 struct {
	B1 [10000]int
}

var r6Pool = sync.Pool{
	New: func() interface{} {
		return new(structR6)
	},
}

// 从堆申请空间
func standardHeap() {
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		var sr6 = new(structR6)
		sr6.B1[i] = 1
	}
	fmt.Println("standardHeap Used:", time.Since(startTime))
}

// 从栈申请空间
func standardStack() {
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		var sr6 structR6
		sr6.B1[i] = 2
	}
	fmt.Println("standardStack Used:", time.Since(startTime))
}

func usePool() {
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		sr6 := r6Pool.Get().(*structR6)
		sr6.B1[i] = 0
		r6Pool.Put(sr6)
	}
	fmt.Println("pool Used:", time.Since(startTime))
}

func main() {
	standardHeap()
	standardStack()
	usePool()

	//aa := r6Pool.Get()
	//fmt.Println(aa)
	//bb := r6Pool.Get()
	//fmt.Println(bb)
	//cc := r6Pool.Get()
	//fmt.Println(cc)

}
