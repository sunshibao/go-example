package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB
var realNum int

func GetDatabase() *gorm.DB {
	return DB
}

func main() {
	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := gorm.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	shell()

	return

	//shell(mongodb, database, skip, limit)

}

func shell() (err error) {

	sql := "select permission from oz_apk_permission"
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return err
	}
	mapTemp := make(map[string]int)
	for rows.Next() {
		var permission string
		rows.Scan(&permission)
		slice := strings.Split(permission, ",")
		for _, v := range slice {
			mapTemp[v] = 1
		}
	}
	for k, _ := range mapTemp {
		InsertTempPermission(k)
	}

	return nil
}

var k = 0

func InsertTempPermission(permission string) {
	k++
	fmt.Println(k)
	sql := "insert into temp_permission (permission) values (?)"
	DB.Exec(sql, permission)
}
