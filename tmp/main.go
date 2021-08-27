package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	s := 0
	for {
		if s >= 100 {
			break
		}
		select {
		case ch <- 0:
			s++
		case ch <- 1:
			s++
		}
		i := <-ch
		fmt.Println("Value received:", i, s)
	}

}
