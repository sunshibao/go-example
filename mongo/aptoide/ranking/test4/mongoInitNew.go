package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type MongoAppInfo struct {
	PackageName string `bson:"package_name"`
	MaxInstalls int64  `bson:"maxInstalls"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	database := "test"
	mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@43.131.69.147:27017,43.131.92.130:27017,162.62.196.12:27017/?replicaSet=market")
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 0
	limit := 100
	s := 0
	var err2 error
	for {
		if err2 == nil && skip < 400000 {
			skip = 0 + limit*s
			err2 = shell(mongodb, database, skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

//10万
//ObjectId("618d74bf685adbd671ac032d")
func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {
	c := mongodb.DB(database).C("package_name")
	var mongoAppInfos []MongoAppInfo
	c.Find(bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex("618d74bf685adbd671ac032d")}}).Select(bson.M{"package_name": 1, "maxInstalls": 1}).Sort("_id").Skip(skip).Limit(limit).All(&mongoAppInfos)
	if len(mongoAppInfos) <= 0 {
		return err
	}

	for k, v := range mongoAppInfos {
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", v.PackageName, k+skip)
		apkSql := "insert into gp_apk (package_name,gp_down_num)values(?,?)"
		DB.Exec(apkSql, v.PackageName, v.MaxInstalls)
	}
	return nil
}
