package main

import (
	"fmt"
	"sync"

	"github.com/robfig/cron"
)

func main() {
	spec := "0 */1 * * * *" // 每分钟
	c := cron.New()
	c.AddFunc(spec, myFunc)
	c.Start()
	select {}
}

func myFunc() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 12; i++ {
		wg.Add(1)
		minId := i * 20000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
}

func start(id int) {
	fmt.Println(id)
}
