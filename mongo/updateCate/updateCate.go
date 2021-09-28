package main

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type MongoAppInfo struct {
	PackageName         string   `bson:"package_name"`
	PackageRlang        string   `bson:"packageRlang"`
	CheckVersionCode    int      `bson:"check_version_code"`
	CheckVersionName    string   `bson:"check_version_name"`
	Description         string   `bson:"description"`
	HeaderImage         string   `bson:"headerImage"`
	Icon                string   `bson:"icon"`
	Installs            string   `bson:"installs"`
	GenreId             string   `bson:"genreId"`
	CheckUserPermission []string `bson:"check_user_permission"`
	PrivacyPolicy       string   `bson:"privacyPolicy"`
	CheckSize           int64    `bson:"check_size"`
	Title               string   `bson:"title"`
	Url                 string   `bson:"url"`
	RecentChanges       string   `bson:"recentChanges"`
	Developer           string   `bson:"developer"`
	Screenshots         []string `bson:"screenshots"`
	Suffix              string   `bson:"suffix"`
}

var DB *gorm.DB

func GetDatabase() *gorm.DB {
	return DB
}

func main() {

	uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 20000 {
			skip = limit * s
			err2 = shell(skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

func shell(start, limit int) (err error) {
	apk_id := 0
	png_icon_id := 0
	package_name := ""
	apkSql := "select apk_id,png_icon_id, package_name from oz_apk_desc t1 limit ?,?"
	DB.Raw(apkSql, start, limit).Row().Scan(&apk_id, &png_icon_id, &package_name)
	fmt.Println("---:", package_name, "---num:", start)
	if apk_id == 0 {
		return err
	}

	apkSql2 := "UPDATE oz_apk_image t1 LEFT JOIN oz_image  t2 ON t1.image_id = t2.image_id SET t2.hd_image_url = REPLACE (t2.hd_image_url,?,?),t2.nhd_image_url = REPLACE (t2.hd_image_url,?,?) WHERE t1.apk_id = ?"
	DB.Exec(apkSql2, "http://18.177.149.123:8001/new/package_name222", "http://18.177.149.123:8001/new/"+package_name, "http://18.177.149.123:8001/new/package_name222", "http://18.177.149.123:8001/new/"+package_name, apk_id)

	apkSql3 := "UPDATE oz_image t2 SET t2.hd_image_url = REPLACE (t2.hd_image_url,?,?),t2.nhd_image_url = REPLACE (t2.hd_image_url,?,?) WHERE t2.image_id = ?"
	DB.Exec(apkSql3, "http://18.177.149.123:8001/new/package_name222", "http://18.177.149.123:8001/new/"+package_name, "http://18.177.149.123:8001/new/package_name222", "http://18.177.149.123:8001/new/"+package_name, png_icon_id)

	return nil
}
