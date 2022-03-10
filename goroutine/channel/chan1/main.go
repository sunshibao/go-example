package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func ch() {
	var ch = make(chan int)
	go func(ch chan int) {
		// ch 没有设置长度所以是阻塞的，会逐个发送
		ch <- 1
		ch <- 2
		ch <- 3
		fmt.Println("send finished")
	}(ch)

	for {
		select {
		case i := <-ch:
			fmt.Println("receive", i)
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

func chLimit() {
	var ch = make(chan int)
	// Tip: 当参数设置为chan<- 和<-chan，可以有效的防止误用发送、接受。
	go func(ch chan<- int) {
		// ch 没有设置长度所以是阻塞的，会逐个发送
		ch <- 1
		ch <- 2
		ch <- 3
		fmt.Println("send finished")
	}(ch)

	for {
		select {
		case i := <-ch:
			fmt.Println("receive", i)
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

// 最标准的，在channel使用完要记得关闭
func chClose() {
	var ch = make(chan int)
	// Tip: 当参数设置为chan<- 和<-chan，可以有效的防止误用发送、接受。
	go func(ch chan<- int) {
		// ch 没有设置长度所以是阻塞的，会逐个发送
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		fmt.Println("send finished")
	}(ch)

	for {
		select {
		case i, ok := <-ch:
			if ok {
				fmt.Println("receive", i)
			} else {
				fmt.Println("channel close")
				os.Exit(1)
			}
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

func chCloseErr() {
	var ch = make(chan int)
	// Tip: 当参数设置为chan<- 和<-chan，可以有效的防止误用发送、接受。
	go func(ch chan<- int) {
		// ch 没有设置长度所以是阻塞的，会逐个发送
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		fmt.Println("send finished")
	}(ch)

	for {
		select {
		case i := <-ch: // 如果这里不接受第二个参数ok,不判断的话，因为channel已经关闭了，所以<-ch一直有值0.永不会停止
			fmt.Println("receive", i)
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

// 用channel做任务调度
func chTask() {
	var doneCh = make(chan struct{})
	var errCh = make(chan error)

	go func(doneCh chan<- struct{}, errCh chan error) {
		if time.Now().Unix()%2 == 0 {
			doneCh <- struct{}{}
		} else {
			errCh <- errors.New("unix time is an odd")
		}
	}(doneCh, errCh)

	for {
		select {
		// 这是一个常见的goroutine 处理模式，在这里监听channel 结果和错误
		case <-doneCh:
			fmt.Println("done")
		case err := <-errCh:
			fmt.Println("get an err", err)
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

func chBuffer() {
	var ch = make(chan int, 3)
	go func(ch chan int) {
		// ch 由于设置了长度，所以这里是不会阻塞的。相当于一个消息队列
		ch <- 1
		ch <- 2
		ch <- 3
		fmt.Println("send finished")
	}(ch)

	for {
		select {
		case i := <-ch:
			fmt.Println("receive", i)
		case <-time.After(time.Second): // 读不到数值了，超时
			fmt.Println("time out")
			os.Exit(1)
		}
	}
}

func chBufferRange() {
	var ch = make(chan int, 3)
	go func(ch chan int) {
		// ch 由于设置了长度，所以这里是不会阻塞的。相当于一个消息队列
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		fmt.Println("send finished")
	}(ch)

	for i := range ch { // 非阻塞的chan 可以用range去接收结果
		fmt.Println(i)
	}
}

func main() {
	//ch()
	//chLimit()
	//chClose()
	//chCloseErr()
	//chBuffer()
	chBufferRange()

}
