package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type LangAppInfo struct {
	DescriptionHtml string `bson:"descriptionHtml"`
	Icon            string `bson:"icon"`
}

type MongoAppInfo struct {
	Package    string                 `bson:"package"`
	UpdateTime time.Time              `bson:"updateTime"`
	Language   map[string]LangAppInfo `bson:"language"`
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
		if err2 == nil && skip < 49999 {
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

//ObjectId("6177e34b275289742a6cf720")
func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {
	c := mongodb.DB(database).C("info")
	var mongoAppInfos []MongoAppInfo
	c.Find(bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex("6177e34b275289742a6cf720")}}).Select(bson.M{"package": 1, "updateTime": 1, "language": 1}).Skip(skip).Limit(limit).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		return err
	}

	for k, s := range mongoAppInfos {
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", s.Package, k+skip)

		// 检查包是否存在
		var apkId int64
		sql := "select apk_id from oz_apk where package_name = ?"
		DB.Get(&apkId, sql, s.Package)

		if apkId > 0 {

			err = insertDB(s, "zh_CN", apkId) //简体中文
			if err != nil {
				continue
			}
			err = insertDB(s, "en_US", apkId) //英语
			if err != nil {
				continue
			}
			err = insertDB(s, "ru_RU", apkId) //俄语
			if err != nil {
				continue
			}
			err = insertDB(s, "be_BY", apkId) //克兰语
			if err != nil {
				continue
			}
			err = insertDB(s, "uk_UA", apkId) //白俄罗斯语
			if err != nil {
				continue
			}
			err = insertDB(s, "kk_KZ", apkId) //哈萨克语
			if err != nil {
				continue
			}
			err = insertDB(s, "ka_GE", apkId) //格鲁吉亚语
			if err != nil {
				continue
			}
			err = insertDB(s, "uz_UZ", apkId) //乌兹别克语
			if err != nil {
				continue
			}

		}
	}
	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64) error {
	v := s.Language[lang]
	var cstZone = time.FixedZone("CST", 0)
	newIcon := "http://18.177.149.123:8001/pic/img_" + s.UpdateTime.In(cstZone).Format("2006-01-02") + "/" + s.Package

	if lang == "zh_CN" {
		lang = "zh_cn"
	} else if lang == "en_US" {
		lang = "en"
	} else if lang == "ru_RU" {
		lang = "ru"
	} else if lang == "be_BY" {
		lang = "be"
	} else if lang == "uk_UA" {
		lang = "uk"
	} else if lang == "kk_KZ" {
		lang = "kk"
	} else if lang == "ka_GE" {
		lang = "ka"
	} else if lang == "uz_UZ" {
		lang = "uz"
	} else {
		lang = "zh_cn"
	}

	statusImg := 1

	v.Icon = strings.Replace(v.Icon, "https://play-lh.googleusercontent.com", newIcon, 1)

	//写入oz_image表
	iconSql3 := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	iconResult, err := DB.Exec(iconSql3, s.Package+"_icon", v.Icon, v.Icon, 50, lang, statusImg)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}
	iconId, _ := iconResult.LastInsertId()

	//写入oz_image表
	imgSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	_, err = DB.Exec(imgSql, s.Package+"_imgUrl", v.Icon, v.Icon, 50, lang, statusImg)

	if err != nil {
		log.Printf("oz_image 2 err:%v", err)
		return err
	}

	apkDescSql := "update oz_apk_desc  set description = ?, png_icon_id = ?,jpg_icon_id = ? where apk_id = ? and language = ?"
	DB.Exec(apkDescSql, v.DescriptionHtml, iconId, iconId, apkId, lang)

	return nil
}
