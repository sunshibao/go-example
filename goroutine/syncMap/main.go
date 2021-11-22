package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

/**
 	首先map 不能直接使用
	错误示范：
	```
	var m map[string]string
	m["result"] = "result"
	```

	在使用时一定要先make
	正确示范：
	```
	m := make(map[string]string)
	m["result"] = "result"
	```
*/

func syncMap() {
	var m sync.Map
	m.Store("a", 1)
	fmt.Println(m.Load("a"))
	// loadOrStore 有就加载，没有就保存值
	fmt.Println(m.LoadOrStore("a", 1))

	m.Delete("a")
	fmt.Println(m.LoadOrStore("a", 2))
	fmt.Println(m.LoadOrStore("b", 3))
	fmt.Println(m.LoadOrStore("c", 4))

	// 遍历
	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true // false 只打印一个。true 全部打印
	})

}

type M struct {
	Map  map[string]string
	lock sync.RWMutex // 加锁
}

// Set ...
func (m *M) Set(key, value string) {
	m.Map[key] = value
}

// Set ...
func (m *M) Set2(key, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Map[key] = value
}

// Get ...
func (m *M) Get(key string) string {
	return m.Map[key]
}
func TestMap() {
	c := M{Map: make(map[string]string)}
	wg := sync.WaitGroup{}
	for i := 0; i < 21; i++ {
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

func TestMap2() {
	c := M{Map: make(map[string]string)}
	wg := sync.WaitGroup{}
	for i := 0; i < 21; i++ {
		wg.Add(1)
		go func(n int) {
			k, v := strconv.Itoa(n), strconv.Itoa(n)
			c.Set2(k, v)
			log.Printf("k=:%v,v:=%v\n", k, c.Get(k))
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("ok finished.")
}

// TestMap  ...
func TestSyncMap() {
	var m sync.Map
	wg := sync.WaitGroup{}
	for i := 0; i < 21; i++ {
		wg.Add(1)
		go func(n int) {
			k, v := strconv.Itoa(n), strconv.Itoa(n)
			m.Store(k, v)
			//load, _ := m.Load(k)                     // 第一种读sync.map方式
			//log.Printf("k=:%v,v:=%v\n", k, load)
			wg.Done()
		}(i)
	}
	wg.Wait()

	m.Range(func(key, value interface{}) bool { // 第二种读sync.map方式
		log.Printf("range k:%v,v=%v\n", key, value)
		return true
	})
}

func main() {
	//TestMap()     // test13 因为map并不是并发安全的，给map写入的时候发生 fatal error: concurrent map writes
	//TestMap2()    // test1 通过加锁方式，解决 map 并发安全
	TestSyncMap() // test3  通过sync包的sync.Map 解决 map 并发安全
}
