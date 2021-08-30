package outGoroutine

import (
	"runtime"
	"testing"
	"time"
)

func test(t *testing.T, f func(chan bool)) {
	t.Helper()
	for i := 0; i < 1000; i++ {
		timeout(f)
	}
	time.Sleep(time.Second * 2)
	t.Log("goroutine num:", runtime.NumGoroutine())
}

func TestBadTimeout(t *testing.T) { test(t, doBadthing) }

func TestBufferTimeout(t *testing.T) {
	for i := 0; i < 1000; i++ {
		timeoutWithBuffer(doBadthing)
	}
	time.Sleep(time.Second * 3)
	t.Log("goroutine num:", runtime.NumGoroutine())
}

func TestGoodTimeout(t *testing.T) { test(t, doGoodthing) }

func Test2phasesTimeout(t *testing.T) {
	for i := 0; i < 1000; i++ {
		timeoutFirstPhase()
	}
	time.Sleep(time.Second * 3)
	t.Log("goroutine num:", runtime.NumGoroutine())
}
