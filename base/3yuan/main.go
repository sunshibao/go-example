package main

import (
	"github.com/sunshibao/go-utils/base"
)

func main()  {
	a, b := "你好", "吗"
	max := base.If(a=="吗", a, b).(string)
	println(max)
}

