package otherGoroutime

import (
	"runtime"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	t.Log(runtime.NumGoroutine())
	sendTasks()
	time.Sleep(time.Second)
	t.Log(runtime.NumGoroutine())
}
