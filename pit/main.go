package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

//1. 可变参数是空接口类型
// 当参数的可变参数是空接口类型时，传入空接口的切片时需要注意参数展开的问题。
func main1() {
	var a = []interface{}{1, 2, 3}

	fmt.Println(a)
	fmt.Println(a...)
}

//2. 数组是值传递
// 在函数调用参数中，数组是值传递，无法通过修改数组类型的参数返回结果。  必要时需要使用切片。
func main2() {
	x := [3]int{1, 2, 3}

	func(arr [3]int) {
		arr[0] = 7
		fmt.Println(arr)
	}(x)

	fmt.Println(x)
}

//3. map遍历是顺序不固定
// map是一种hash表实现，每次遍历的顺序都可能不一样。
func main3() {
	m := map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
	}

	for k, v := range m {
		println(k, v)
	}
}

//3. recover必须在defer函数中运行

// recover捕获的是祖父级调用时的异常，直接调用时无效：
//func main() {
//	recover()
//	panic(1)
//}

//直接defer调用也是无效：
//func main() {
//	defer recover()
//	panic(1)
//}

//defer调用时多层嵌套依然无效：
//func main() {
//	defer func() {
//		func() { recover() }()
//	}()
//	panic(1)
//}

//必须在defer函数中直接调用才有效：
func main4() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()
	panic(1)
}

//5. main函数提前退出
// 后台Goroutine无法保证完成任务。 在使用协程时一定要使用sync.waitGroup,保证协程执行完main在退出
func main5() {
	go println("hello")
}

//6. 通过Sleep来回避并发中的问题
// 休眠并不能保证输出完整的字符串：
func main6() {
	go println("hello")
	time.Sleep(time.Second)
}

//类似的还有通过插入调度语句：
func main7() {
	go println("hello")
	runtime.Gosched()
}

//8. 独占CPU导致其它Goroutine饿死
// Goroutine是协作式抢占调度，Goroutine本身不会主动放弃CPU：
func main8() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	for {
	} // 占用CPU
}

//解决的方法是在for循环加入runtime.Gosched()调度函数：
func main9() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	for {
		runtime.Gosched()
	}
}

//或者是通过阻塞的方式避免CPU占用：
func main10() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		os.Exit(0)
	}()

	select{}
}

//11. 不同Goroutine之间不满足顺序一致性内存模型
// 因为在不同的Goroutine，main函数中无法保证能打印出hello, world:
var msg string
var done bool
func setup() {
	msg = "hello, world"
	done = true
}

func main11() {
	go setup()
	for !done {
	}
	println(msg)
}

//解决的办法是用显式同步：
var msg1 string
var done1 = make(chan bool)
func setup1() {
	msg1 = "hello, world"
	done1 <- true
}

//msg的写入是在channel发送之前，所以能保证打印hello, world
func main12() {
	go setup1()
	<-done1
	println(msg1)
}


//13. 闭包错误引用同一个变量
func main13() {
	for i := 0; i < 5; i++ {
		defer func() {
			println(i)
		}()
	}
}
//改进的方法是在每轮迭代中生成一个局部变量：
func main14() {
	for i := 0; i < 5; i++ {
		i := i
		defer func() {
			println(i)
		}()
	}
}
//或者是通过函数参数传入：
func main15() {
	for i := 0; i < 5; i++ {
		defer func(i int) {
			println(i)
		}(i)
	}
}

//16. 在循环内部执行defer语句
//defer在函数退出时才能执行，在for执行defer会导致资源延迟释放：
func main16() {
	for i := 0; i < 5; i++ {
		f, err := os.Open("/Users/sunshibao/Desktop/技术文档/大老板.md")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

//解决的方法可以在for中构造一个局部函数，在局部函数内部执行defer：
func main17() {
	for i := 0; i < 5; i++ {
		func() {
			f, err := os.Open("/Users/sunshibao/Desktop/技术文档/大老板.md")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
		}()
	}
}

//18. 切片会导致整个底层数组被锁定
//切片会导致整个底层数组被锁定，底层数组无法释放内存。如果底层数组较大会对内存产生很大的压力。
func main18() {
	headerMap := make(map[string][]byte)

	for i := 0; i < 5; i++ {
		name := "/path/to/file"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		headerMap[name] = data[:1]
	}

	// do some thing
}

//解决的方法是将结果克隆一份，这样可以释放底层的数组：
func main19() {
	headerMap := make(map[string][]byte)

	for i := 0; i < 5; i++ {
		name := "/path/to/file"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		headerMap[name] = append([]byte{}, data[:1]...)
	}

	// do some thing
}

//20. Goroutine泄露
//Go语言是带内存自动回收的特性，因此内存一般不会泄漏。但是Goroutine确存在泄漏的情况，同时泄漏的Goroutine引用的内存同样无法被回收。

func main20() {
	ch := func() <-chan int {
		ch := make(chan int)
		go func() {
			for i := 0; ; i++ {
				ch <- i
			}
		} ()
		return ch
	}()

	for v := range ch {
		fmt.Println(v)
		if v == 5 {
			break
		}
	}
}

//上面的程序中后台Goroutine向管道输入自然数序列，main函数中输出序列。但是当break跳出for循环的时候，后台Goroutine就处于无法被回收的状态了。
//
//我们可以通过context包来避免这个问题：

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ch := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		go func() {
			for i := 0; ; i++ {
				select {
				case <- ctx.Done():
					return
				case ch <- i:
				}
			}
		} ()
		return ch
	}(ctx)

	for v := range ch {
		fmt.Println(v)
		if v == 5 {
			cancel()
			break
		}
	}
}
//当main函数在break跳出循环时，通过调用cancel()来通知后台Goroutine退出，这样就避免了Goroutine的泄漏。