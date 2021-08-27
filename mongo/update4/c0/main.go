package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type UpdateAppInfo struct {
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
	Version             string   `bson:"version"`
	RecentChanges       string   `bson:"recentChanges"`
	Developer           string   `bson:"developer"`
	Screenshots         []string `bson:"screenshots"`
	Suffix              string   `bson:"suffix"`
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

	//skip := 50000
	//limit := 1
	//s := 0
	//var err2 error
	//for {
	//	if err2 == nil {
	//		skip = 50000 + limit*s
	//		err2 = shell(mongodb, database, skip, limit)
	//		s++
	//	} else {
	//		break
	//	}
	//}
	//return

	shell(mongodb, database)

}

//com.MattGames.SixCups
//com.baskosergey.basketroll
//com.fungame.digsandball
func shell(mongodb *mgo.Session, database string) (err error) {
	c := mongodb.DB(database).C("package")
	var mongoAppInfos []UpdateAppInfo
	t1 := time.Now()
	arr := []string{"com.MattGames.Six", "com.baskosergey.basket", "com.fungame.digsand"}
	c.Find(bson.M{"package_name": bson.M{"$in": bson.M{"$regex": arr}}}).All(&mongoAppInfos)
	elapsed := time.Since(t1)

	fmt.Println(elapsed)
	for k, v := range mongoAppInfos {
		log.Println(k, "---", v.PackageName)
	}
	return nil
}
