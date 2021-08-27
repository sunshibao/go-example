package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UpdateAppInfo struct {
	PackageName         string   `bson:"package_name"`
	CheckVersionCode    int      `bson:"check_version_code"`
	CheckVersionName    string   `bson:"check_version_name"`
	CheckUserPermission []string `bson:"check_user_permission"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	host := "18.177.149.123:27017"
	database := "test"
	mongodb, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 44000
	limit := 1
	s := 0
	var err2 error
	for {
		if err2 == nil && skip <= 45000 {
			skip = 44000 + limit*s
			err2 = shell(mongodb, database, skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {
	c := mongodb.DB(database).C("package")
	var mongoAppInfos []UpdateAppInfo
	c.Find(bson.M{"pic": 3, "apk": 1, "re": bson.M{"$in": []int{1, 2}}, "suffix": bson.M{"$exists": true}, "appId": bson.M{"$exists": true}, "check_size": bson.M{"$exists": true}}).Skip(skip).Limit(limit).All(&mongoAppInfos)

	for _, v := range mongoAppInfos {
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", v.PackageName, skip)
		permission := strings.Join(v.CheckUserPermission, ",")
		apkSql := "UPDATE oz_apk oa LEFT JOIN oz_apk_desc oad ON oa.apk_id=oad.apk_id SET oa.`version_code`= ? ,oa.`version_name`= ? ,oad.app_permission_desc= ?  WHERE oa.package_name= ? "
		_, err := DB.Exec(apkSql, v.CheckVersionCode, v.CheckVersionName, permission, v.PackageName)
		if err != nil {
			log.Printf("apk err:%v\n", err)
			continue
		}
	}
	return nil
}
