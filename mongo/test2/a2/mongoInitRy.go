package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func GetDatabase() *gorm.DB {
	return DB
}

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := gorm.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 0
	limit := 1
	s := 0
	var err2 error
	for {
		if err2 == nil && skip < 6000 {
			skip = 0 + limit*s
			err2 = shell(skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

type AAA struct {
	ApkId       int    `json:"apk_id"`
	PackageName string `json:"package_name"`
	DownloadNum int    `json:"download_num"`
	ApkResType  string `json:"apk_res_type"`
}

func shell(skip, limit int) (err error) {
	// 检查包是否存在
	var apPackage string
	var name string
	sql := "SELECT name,package FROM ws77 where billing = 1  limit ?,?"
	DB.Raw(sql, skip, limit).Row().Scan(&name, &apPackage)

	if name != "" {
		var gpPackage string
		sql2 := "select package_name from 222_oz_apk where apk_name = ? ;"
		DB.Raw(sql2, name).Row().Scan(&gpPackage)

		if gpPackage != "" && gpPackage != apPackage {
			sql3 := "update ws77 set temp_billing =1 where name = ?"
			DB.Exec(sql3, name)
		}
	}
	log.Printf("PackageName:%s-----num:%d\n", apPackage, skip)

	return nil
}
