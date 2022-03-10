package main

import "fmt"

/*
平铺设计模式，适合单业务，业务比较少的情况
*/

//我们要写一个类,Banker银行业务员
type Banker struct {
}

//存款业务
func (this *Banker) Save() {
	fmt.Println("进行了 存款业务...")
}

//转账业务
func (this *Banker) Transfer() {
	fmt.Println("进行了 转账业务...")
}

//支付业务
func (this *Banker) Pay() {
	fmt.Println("进行了 支付业务...")
}

func main() {
	banker := &Banker{}

	banker.Save()
	banker.Transfer()
	banker.Pay()
}
