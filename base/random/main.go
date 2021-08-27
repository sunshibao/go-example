package main

import (
	"fmt"
	"github.com/sunshibao/go-utils/base"
	"math/rand"
)

func main() {
	//1. 基本随机数,但是数值不会变。
	a := rand.Int()
	b := rand.Intn(100) //生成0-99之间的随机数
	fmt.Println(a)
	fmt.Println(b)

	//2. 生成可变随机数, 将时间戳设置成种子数.
	//rand.Seed(time.Now().UnixNano())
	//生成10个0-99之间的随机数
	for i := 0; i < 10; i++ {
		//fmt.Println(rand.Intn(100))
	}

	//3. 生成[15，88]之间的随机数,括号左包含右不包含 数值也不会变，必须设置种子才会变。
	n := rand.Intn(73) + 15 //(88-15 )+15
	fmt.Println(n)

	// 4. 第三方包 生成01的
	for i := range base.Random(10) {
		fmt.Println(i)
	}
}
