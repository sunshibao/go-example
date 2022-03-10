package main

import (
	"context"
	"fmt"
	"time"
)

// Tip 通过cannel 主动关闭
func ctxCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-time.After(time.Microsecond * 100):
			fmt.Println("time out")
		}

	}(ctx)
	cancel()
}

// 通过超时自动触发
func ctxTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*10)

	//主动执行cancel 也会让协程收到消息
	defer cancel()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-time.After(time.Microsecond * 100):
			fmt.Println("time out")
		}

	}(ctx)
	time.Sleep(time.Second)
}

// 通过截止时间自动触发
func ctxDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Microsecond*10))

	//主动执行cancel 也会让协程收到消息
	defer cancel()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-time.After(time.Microsecond * 100):
			fmt.Println("time out")
		}

	}(ctx)
	time.Sleep(time.Second)
}

// 用key/value 传递参数
func ctxValue() {
	ctx := context.WithValue(context.Background(), "user", "sunshibao")
	go func(ctx context.Context) {
		v, ok := ctx.Value("user ").(string)
		if ok {
			fmt.Println("pass user value", v)
		}

	}(ctx)
	time.Sleep(time.Second)
}

func main() {
	//ctxCancel()
	//ctxTimeout()
	//ctxDeadline()
	ctxValue()
}
