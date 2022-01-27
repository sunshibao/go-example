package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *gorm.DB

func main() {
	//for i := 0; i <= 9; i++ {
	//	minId := i * 10000
	//	go func(id int) {
	//		start(id)
	//	}(minId)
	//}
	start(0)
}

func start(minId int) {
	//建立连接
	uri := "root:Droi*#2021@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 100808 {
			skip = 0 + limit*s
			err2 = GetApkList(minId, skip, limit)
			s++
		} else {
			break
		}
	}
	return
}

func GetApkList(id, skip, limit int) (err error) {
	sql1 := "select package,developer_privacy from ws80_detail where id >? order by id limit ?,? "
	rows, err := DB.Raw(sql1, id, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var packageName string
			var developerPrivacy string
			err := rows.Scan(&packageName, &developerPrivacy)
			if err != nil {
				continue
			}
			fmt.Println(packageName, "--------", skip)
			if packageName != "" {
				UploadCos(packageName, developerPrivacy)
			}

		}
	}
	return nil
}

func UploadCos(packageName, developerPrivacy string) {
	sql := `update oz_apk_desc set app_permission_url = ? where package_name = ? and apk_id>438676`
	DB.Exec(sql, developerPrivacy, packageName)
}
