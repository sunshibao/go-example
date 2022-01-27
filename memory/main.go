package main

import (
	"fmt"
	"unsafe"
)

var aa *int
var bb interface{}

func main() {
	part1 := Part1{}
	Part2 := Part2{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
	fmt.Printf("part2 size: %d, align: %d\n", unsafe.Sizeof(Part2), unsafe.Alignof(Part2))

}

type Part1 struct {
	A int8
	j int16
	y int8
	k int32
	g int8
	B int64
}

type Part2 struct {
	a int32 //1  ---3
	c int8
	b int32
	//
}

// {1,1,1,1,1,1,0,0}
// {0,0,0,0,0,0,0,0}
// {1,0,0,0,0,0,0,0}
