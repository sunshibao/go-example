package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type LangAppInfo struct {
	Icon    string   `bson:"icon"`
	ImgList []string `bson:"imgList"`
}

type MongoAppInfo struct {
	Package    string                 `bson:"package"`
	UpdateTime time.Time              `bson:"updateTime1"`
	Language   map[string]LangAppInfo `bson:"language"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Open("mysql", uri)
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
		if err2 == nil && skip < 1 {
			skip = 0 + limit*s
			err2 = shell()
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

//ObjectId("6177e34b275289742a6cf720")
func shell() (err error) {
	//写入oz_image表
	iconSql3 := "delete from oz_image where image_id>3844965 limit 10;"
	DB.Exec(iconSql3)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}

	iconSql4 := "delete from oz_apk_image where apk_id>31196 limit 10;"
	DB.Exec(iconSql4)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}
	return nil
}
