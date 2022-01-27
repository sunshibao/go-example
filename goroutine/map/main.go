package main

import (
	"log"
	"strconv"
	"sync"
)

func main() {
	c := M{Map: make(map[string]string)}
	wg := sync.WaitGroup{}
	for i := 0; i < 210; i++ {
		wg.Add(1)
		go func(n int) {
			k, v := strconv.Itoa(n), strconv.Itoa(n)
			c.Set(k, v)
			log.Printf("k=:%v,v:=%v\n", k, c.Get(k))
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("ok finished.")
}

type M struct {
	Map  map[string]string
	lock sync.RWMutex
}

// Set ...
func (m *M) Set(key, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Map[key] = value
}

// Get ...
func (m *M) Get(key string) string {
	return m.Map[key]
}
