package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type BadgeForLegacyRating struct {
	Major string
}
type Annotations struct {
	BadgeForLegacyRating BadgeForLegacyRating `bson:"badgeForLegacyRating"`
}

type LangAppInfo struct {
	Annotations Annotations `bson:"annotations"`
}

type MongoAppInfo struct {
	Package  string                 `bson:"package"`
	Language map[string]LangAppInfo `bson:"language"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	host := "18.177.149.123:27017"
	database := "package"
	mongodb, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_clone?charset=utf8mb4&parseTime=True&loc=Local"
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
		if err2 == nil {
			skip = limit * s
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
	c := mongodb.DB(database).C("info")
	var mongoAppInfos []MongoAppInfo
	c.Find(nil).Skip(skip).Limit(limit).All(&mongoAppInfos)
	if len(mongoAppInfos) <= 0 {
		return err
	}
	for k, s := range mongoAppInfos {

		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", s.Package, k+skip)
		major := s.Language["zh_CN"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "zh_cn", major)

		major = s.Language["en_US"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "en", major)

		major = s.Language["ru_RU"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "ru", major)

		major = s.Language["be_BY"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "be", major)

		major = s.Language["uk_UA"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "uk", major)

		major = s.Language["ka_GE"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "ka", major)

		major = s.Language["kk_KZ"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "kk", major)

		major = s.Language["uz_UZ"].Annotations.BadgeForLegacyRating.Major
		langShell(s.Package, "uz", major)
	}
	return nil
}

func langShell(packageName, lang, major string) {
	// 检查包是否存在
	var apkId int64
	sql := "select apk_id from oz_apk_desc where package_name = ? and language = ?"
	DB.Get(&apkId, sql, packageName, lang)
	if apkId > 0 {
		apkSql := "update oz_apk_desc set age_limit = ? where package_name = ? and language = ?"
		DB.Exec(apkSql, major, packageName, lang)
	}
}
