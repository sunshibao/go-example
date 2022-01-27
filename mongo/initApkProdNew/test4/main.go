package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type BadgeForLegacyRating struct {
	Major string
}

type Ds struct {
	Low int `json:"low"`
}

type Dc struct {
	Low int `json:"low"`
}

type AppInfo struct {
	DeveloperName     string   `bson:"developerName"`
	DeveloperEmail    string   `bson:"developerEmail"`
	DeveloperWebsite  string   `Bson:"developerWebsite"`
	VersionCode       int      `bson:"versionCode"`
	VersionString     string   `bson:"versionString"`
	Permission        []string `bson:"permission"`
	PackageName       string   `bson:"packageName"`
	AppType           string   `bson:"appType"`
	CategoryName      string   `bson:"categoryName"`
	RecentChangesHtml string   `bson:"recentChangesHtml"`
	InfoDownloadSize  Ds       `bson:"infoDownloadSize"`
	DownloadCount     Dc       `bson:"downloadCount"`
	InstallNotes      string   `bson:"installNotes"`
	InAppProduct      string   `bson:"inAppProduct"`
}

type Annotations struct {
	PrivacyPolicyUrl     string               `bson:"privacyPolicyUrl"`
	BadgeForLegacyRating BadgeForLegacyRating `bson:"badgeForLegacyRating"`
}

type Detail struct {
	AppDetails AppInfo `bson:"appDetails"`
}

type LangAppInfo struct {
	ID                     string      `bson:"id"`
	Type                   int         `bson:"type"`
	CategoryId             string      `bson:"CategoryId"`
	Title                  string      `bson:"title"`
	Creator                string      `bson:"creator"`
	DescriptionHtml        string      `bson:"descriptionHtml"`
	Icon                   string      `bson:"icon"`
	ImgList                []string    `bson:"imgList"`
	Detail                 Detail      `bson:"details"`
	ShareUrl               string      `bson:"shareUrl"`
	Annotations            Annotations `bson:"annotations"`
	DromotionalDescription string      `bson:"promotionalDescription"`
}

type MongoAppInfo struct {
	Package    string                 `bson:"package"`
	Language   map[string]LangAppInfo `bson:"language"`
	UpdateTime time.Time              `bson:"updateTime"`
	Status     int                    `bson:"status"`
}

var DB *sqlx.DB
var realNum int

func GetDatabase() *sqlx.DB {
	return DB
}

func main() {

	database := "package"
	mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@43.131.69.147:27017,43.131.92.130:27017,162.62.196.12:27017/?replicaSet=market")
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

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
		if err2 == nil && skip < 485 {
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

//614ade639a408f2d02455321
//com.enax.zombieevilkill3
//com.nexamuse.cakedesigns
func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {
	c := mongodb.DB(database).C("info")
	var mongoAppInfos []MongoAppInfo

	c.Find(bson.M{"package": "com.nexamuse.cakedesigns"}).Sort("_id").Skip(skip).Limit(limit).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		return err
	}

	for k, s := range mongoAppInfos {
		// 检查包是否存在
		var apkId int64
		sql := "select apk_id from oz_apk where package_name = ?"
		DB.Get(&apkId, sql, s.Package)

		log.Printf("PackageName1111111:%s-----num:%d\n", s.Package, k+skip)

		err = insertDB(s, "ru_RU", 2) //俄语
		if err != nil {
			continue
		}

	}
	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64) error {
	h, _ := time.ParseDuration("-1h")

	newIcon := "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/app_img/img_" + s.UpdateTime.Add(8*h).Format("2006-01-02") + "/" + s.Package
	fmt.Println(s.UpdateTime)
	fmt.Println(s.UpdateTime.Add(8 * h))
	fmt.Println(newIcon)
	return nil
}
