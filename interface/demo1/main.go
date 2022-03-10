package main

import "fmt"

/**
张三，李四。。。分别开奔驰，宝马。
*/

//抽象层
type Car interface {
	Run()
}

type Driver interface {
	Drive(car Car)
}

//实现层
type Benchi struct {
}

type Baoma struct {
}

type Zhangsan struct {
}

type Lisi struct {
}

func (baoma *Baoma) Run() {
	fmt.Println("baoma is running...")
}

func (benchi *Benchi) Run() {
	fmt.Println("benchi is running...")
}

func (lisi *Lisi) Drive(car Car) {
	fmt.Println("lisi driver")
	car.Run()
}

func (zhangsan *Zhangsan) Drive(car Car) {
	fmt.Println("zhangsan driver")
	car.Run()
}

//业务逻辑层
func main() {
	var baoma Car
	baoma = &Baoma{}

	var zhangsan Driver
	zhangsan = &Zhangsan{}

	zhangsan.Drive(baoma)

	var benchi Car
	benchi = &Benchi{}

	var lisi Driver
	lisi = &Lisi{}

	lisi.Drive(benchi)
}
