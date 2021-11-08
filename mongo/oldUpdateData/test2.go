package main

import (
	"fmt"
	"strings"

	"github.com/sunshibao/go-utils/util/gconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type MongoAppInfo struct {
	PackageName string `bson:"package_name"`
	Installs    string `bson:"installs"`
	Size        string `Bson:"size"`
	Version     string `bson:"version"`
	Url         string `bson:"url"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	host := "18.177.149.123:27017"
	database := "test13"
	mongodb, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 300000
	limit := 1000
	s := 0
	var err2 error
	for {
		if err2 == nil && skip < 400000 {
			skip = 300000 + limit*s
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
	c := mongodb.DB(database).C("package_name")
	var mongoAppInfos []MongoAppInfo
	c.Find(nil).Skip(skip).Limit(limit).All(&mongoAppInfos)
	if len(mongoAppInfos) <= 0 {
		return err
	}
	for _, v := range mongoAppInfos {
		installs := strings.Replace(v.Installs, " ", "", -1)
		installs = strings.Replace(installs, ",", "", -1)
		installs = strings.Replace(installs, "+", "", -1)
		code := strings.Replace(v.Version, ".", "", -1)
		fileSizeString := strings.Replace(v.Size, "M", "", -1)
		size := gconv.Int(fileSizeString) * 1024 * 1024
		apkSql := "update oz_apk set  gp_down_url = ?,file_size = ?,version_code = ?,version_name = ?,gp_down_url = ? where package_name = ?"

		DB.Exec(apkSql, installs, size, code, v.Version, v.Url, v.PackageName)
		fmt.Println(v.PackageName, skip)
	}
	return nil
}
