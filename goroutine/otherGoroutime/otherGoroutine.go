package otherGoroutime

import (
	"fmt"
	"time"
)

func do(taskCh chan int) {
	for {
		select {
		case _, beforeClosed := <-taskCh:
			if !beforeClosed { //要判断是否关闭，如果关闭直接退出协程
				fmt.Println("taskCh has been closed")
				return
			}
			time.Sleep(time.Millisecond)
			//fmt.Printf("task %d is done\n", t)
		}
	}
}

func sendTasks() {
	taskCh := make(chan int, 10)
	go do(taskCh)
	for i := 0; i < 1000; i++ {
		taskCh <- i
	}
	close(taskCh) //一定要及时关闭
}
