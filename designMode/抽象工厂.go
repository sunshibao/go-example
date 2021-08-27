package designMode

import (
	"fmt"
)

// 创建型-抽象工厂模式

type Girl2 interface {
	weight()
}

// 中国胖女孩
type FatGirl2 struct {
}

func (FatGirl2) weight() {
	fmt.Println("chinese girl weight: 80kg")
}

// 瘦女孩
type ThinGirl2 struct {
}

func (ThinGirl2) weight() {
	fmt.Println("chinese girl weight: 50kg")
}

type Factory interface {
	CreateGirl(like string) Girl2
}

// 中国工厂
type ChineseGirlFactory struct {
}

func (ChineseGirlFactory) CreateGirl(like string) Girl2 {
	if like == "fat" {
		return &FatGirl2{}
	} else if like == "thin" {
		return &ThinGirl2{}
	}
	return nil
}

// 美国工厂
type AmericanGirlFactory struct {
}

func (AmericanGirlFactory) CreateGirl(like string) Girl2 {
	if like == "fat" {
		return &AmericanFatGirl{}
	} else if like == "thin" {
		return &AmericanThainGirl{}
	}
	return nil
}

// 美国胖女孩
type AmericanFatGirl struct {
}

func (AmericanFatGirl) weight() {
	fmt.Println("American weight: 80kg")
}

// 美国瘦女孩
type AmericanThainGirl struct {
}

func (AmericanThainGirl) weight() {
	fmt.Println("American weight: 50kg")
}

// 工厂提供者
type GirlFactoryStore struct {
	factory Factory
}

func (store *GirlFactoryStore) createGirl(like string) Girl2 {
	return store.factory.CreateGirl(like)
}