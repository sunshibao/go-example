package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func game1() {
	const (
		N         = 10 //10个人
		M         = 10 // 10轮游戏
		MaxNumber = 10
		gameOver  = 0
		gameNext  = 1
	)
	ch := make(chan chan int) //这里定义了一个管道类型的管道
	readyGame := sync.WaitGroup{}
	readyGame.Add(N)
	start := sync.WaitGroup{} //start是用来控制每轮游戏暂停，保证没有人先开下一轮的
	// N个人
	for i := 0; i < N; i++ {
		go func() {
			//每个人执行M轮游戏
			phone := make(chan int)
			for j := 0; j < M; j++ {
				// 告诉裁判我可以开始游戏了
				readyGame.Done()
				ch <- phone //注意这里是将phone给了ch，不是从phone里取值
				phone <- rand.Intn(MaxNumber)
				//裁判是通过ch里面的phone里面的val来判断你有没有出局
				if val := <-phone; val == gameOver {
					break
				}
				//这一轮游戏结束等待
				// 如果没有这一步的话 他会直接进入到写管道 这个时候他可能提前进入下一轮游戏 这样是不对的
				start.Wait()
			}

		}()
	}
	//裁判
	start.Add(1)
	cnt := N
	for i := 0; i < M; i++ {
		num := rand.Intn(MaxNumber) //裁判选数字
		tempCnt := cnt
		for j := 0; j < tempCnt; j++ {
			phone := <-ch
			val := <-phone
			//fmt.Println(val)
			if num == val {
				phone <- gameOver
				cnt--
			} else {
				phone <- gameNext
			}
		}
		if i == M-1 {
			break
		}
		readyGame.Add(cnt) //给剩下的人发通行证
		start.Done()       //将裁判的1变为0
		readyGame.Wait()   //保证下一次开始前选手们都已经选好
		start.Add(1)
	}
	fmt.Println(M, " 轮游戏还有", cnt, "个人")
}

func main() {
	game1()
}