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
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 100000 {
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

//614ade639a408f2d02455321
func shell(skip, limit int) (err error) {
	// 检查包是否存在
	var packageName string
	var isGp int
	sql := "select package,is_gp from ws76 order by id limit ?,?"
	DB.Raw(sql, skip, limit).Row().Scan(&packageName, &isGp)
	if isGp == 1 {
		var apkName string
		sql2 := "select apk_name from 222_oz_apk where package_name = ? ;"
		DB.Raw(sql2, packageName).Row().Scan(&apkName)

		if apkName != "" {
			sql3 := "update ws76 set gp_package_name = ?,gp_apk_name =? where package = ?"
			DB.Exec(sql3, packageName, apkName, packageName)
		}

	}

	log.Printf("PackageName:%s-----num:%d\n", packageName, skip)

	return nil
}
