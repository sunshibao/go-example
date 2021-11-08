package main

import "testing"

// go test13 -bench='Fib$' -benchtime=5s -count=3 .
// -bench 哪个文件 正则匹配   `Fib$` 以Fib结尾。   `Generate` 包含Generate
// -benchtime 多少秒（5s）。。。多少个（30x）
// -count 多少轮
// -benchmem 参数看到内存分配的情况

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

func TestFib(t *testing.T) {
	for n := 0; n < 20; n++ {
		fib(30) // run fib(30) b.N times
	}
}
