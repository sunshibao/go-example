package main

import "fmt"

// Sleeper 接口声明
type Sleeper interface {
	Sleep() // 声明一个Sleep() 方法
}

type Dog struct {
	Name string
}

type Cat struct {
	Name string
}

func (d Dog) Sleep() {
	fmt.Printf("Dog %s is sleeping\n", d.Name)
}

func (c Cat) Sleep() {
	fmt.Printf("Cat %s is sleeping\n", c.Name)
}

func AnimalSleep(s Sleeper) { // 注意参数是一个 interface
	s.Sleep()
}
func main() {
	var s Sleeper
	dog := Dog{Name: "xiaobai"}
	cat := Cat{Name: "hellokitty"}

	s = dog
	AnimalSleep(s)
	s = cat
	AnimalSleep(s)

	// 创建一个Sleeper 切片
	sleepList := []Sleeper{Dog{Name: "xiaobai"}, Cat{Name: "hellokitty"}}
	for _, s := range sleepList {
		s.Sleep()
	}
}
