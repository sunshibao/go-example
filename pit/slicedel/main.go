package main

import (
	"container/list"
	"fmt"
)

func main() {

	l := list.New()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	fmt.Println("original list:")
	prtList(l)

	fmt.Println("deleted list:")

	for e := l.Front(); e != nil; e = e.Next() {
		l.Remove(e)
	}

	prtList(l)
}

func prtList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v ", e.Value)
	}
	fmt.Printf("\n")
}


//由源码中可以看到，当执行l.Remove(e)时，会在内部调用l.remove(e)方法删除元素e，为了避免内存泄漏，会将e.next和e.prev赋值为nil，这就是问题根源。

func main2() {

	l := list.New()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	fmt.Println("original list:")
	prtList(l)

	fmt.Println("deleted list:")
	var next *list.Element
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}

	prtList(l)
}