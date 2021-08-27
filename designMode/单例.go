package designMode

import (
	"sync"
)

// 创建型-单例模式

//使方法只执行一次的对象实现，作用与init函数类似
//init函数是在文件包首次被加载的时候执行，且只执行一次
//sync.Onc是在代码运行中需要的时候执行，且只执行一次

// 单例模式之懒汉模式。
type singleton struct{}

var (
	instance *singleton
	once     sync.Once
)

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

// 单例模式之饿汉模式
type singleton2 struct{}

var (
	instance2 *singleton2
)

func init() {
	instance2 = &singleton2{}
}

func GetInstance2() *singleton2 {
	return instance2
}
