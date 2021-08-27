package main

import (
	"fmt"
	log "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type AppInfo struct {
	PackageName string `json:"package_name"`
	ApkResType  string `json:"apk_res_type"`
}

var DB *gorm.DB

func main() {
	var err error
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	mysqldb, err := gorm.Open("mysql", uri)
	if err != nil {
		log.Error("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	package_name := "com.os.airforce"
	// sql1.查询一条数据 一列或多列   scan
	type ApkData struct {
		ApkResType string
	}
	var apkData ApkData
	sql := `select apk_res_type from oz_apk where package_name = ?`
	err = DB.Raw(sql, package_name).Scan(&apkData).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql1:",apkData)
	// sql2.查询一条数据 一列    pluck
	apkRes := []string{}
	sql1 := `select apk_res_type from oz_apk where package_name = ?`
	err3 := DB.Raw(sql1, package_name).Pluck("apk_res_type", &apkRes).Error
	if err3 != nil {
		log.Error(err)
	}
	fmt.Println("sql2:",apkRes)

	// sql3查询多条
	newPackageName := []string{}
	apkSql := `select package_name from oz_apk limit 10`
	rows, err := DB.Raw(apkSql).Rows()

	for rows.Next() {
		packageName := ""
		err := rows.Scan(&packageName)
		if err != nil {
			continue
		}
		newPackageName = append(newPackageName, packageName)
	}
	fmt.Println("sql3:",newPackageName)

}
