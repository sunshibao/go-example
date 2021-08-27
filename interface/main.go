package main

import "fmt"

type Dog interface {
	Hi()
}

type Dog1 struct {}

func (d1 Dog1) Hi()  {}

type Dog2 struct {}

func (d2 Dog2) Hi()  {}

func dataInterface()  {
	var i interface{} = 1
	fmt.Println(i)
	// 类型定义，除非确定100%成功。尽量用两个返回参数，否则会pannic
	v,ok:= i.(string)
	fmt.Printf(v,ok)
	var d1 Dog1
	var d2 Dog2

	var dList = []Dog{d1,d2}
	for _,v:=range dList{
		v.Hi()
	}
}
func main()  {
	dataInterface()
}