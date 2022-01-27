package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 0)

	go func() {
		ch <- 666
	}()

	go func() {
		x := <-ch
		fmt.Println(x)
	}()
	time.Sleep(2)
	return
}
