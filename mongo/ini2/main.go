package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sunshibao/go-utils/util/gconv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	uri := "root:sun18188@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

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
	c := mongodb.DB(database).C("package")
	var mongoAppInfos []MongoAppInfo
	c.Find(bson.M{"pic": 3, "apk": 1, "genreId": bson.M{"$regex": "GAME"}, "re": bson.M{"$in": []int{1, 2}}, "suffix": bson.M{"$exists": true}, "appId": bson.M{"$exists": true}, "check_size": bson.M{"$exists": true}, "_id": bson.M{"$gt": bson.ObjectIdHex("6119fc57906ddd5816c8a5c0")}}).Skip(skip).Limit(limit).All(&mongoAppInfos)

	for _, v := range mongoAppInfos {
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", v.PackageName, skip)
		newIcon := "http://18.177.149.123:8001/resources/" + v.PackageName + "/icon_"
		newImg := "http://18.177.149.123:8001/resources/" + v.PackageName + "/headerImage_"
		downloadUrl := "http://18.177.149.123:8001/apk/" + v.Suffix

		index := strings.Contains(v.Suffix, ".apk")

		statusApk := 1
		statusApkDesc := 1
		statusImg := 1
		if !index {
			statusApk = -1
			statusApkDesc = 2
			statusImg = 0
		}

		apkSql := "insert into kk_apk (package_name,apk_name,version_code,version_name,download_url,company,file_size,download_num,apk_res_type,status)values(?,?,?,?,?,?,?,?,?,?)"
		appResult, err := DB.Exec(apkSql, v.PackageName, v.Title, v.CheckVersionCode, v.CheckVersionName, downloadUrl, v.Developer, v.CheckSize, gconv.Int64(strings.Replace(v.Installs, "+", "", 1)), strings.ToLower(v.GenreId), statusApk)

		if err != nil {
			log.Printf("oz_apk err:%v", err)
			continue
		}
		apkId, _ := appResult.LastInsertId()

		//处理资源
		img_url := strings.Replace(v.Icon, "https://play-lh.googleusercontent.com/", newImg, 1)
		icon := strings.Replace(v.Icon, "https://play-lh.googleusercontent.com/", newIcon, 1)
		//写入oz_image表
		iconSql := "INSERT INTO kk_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
		iconResult, err := DB.Exec(iconSql, v.PackageName+"_icon", icon, icon, 50, "ru", statusImg)

		if err != nil {
			log.Printf("oz_image 1 err:%v", err)
			continue
		}
		iconId, _ := iconResult.LastInsertId()
		//写入oz_image表
		imgSql := "INSERT INTO kk_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
		_, err = DB.Exec(imgSql, v.PackageName+"_imgUrl", img_url, img_url, 50, "ru", statusImg)

		if err != nil {
			log.Printf("oz_image 2 err:%v", err)
			continue
		}
		permission := strings.Join(v.CheckUserPermission, ",")
		apkDescSql := "insert into kk_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,status)values(?,?,?,?,?,?,?,?,?,?,?)"
		_, err = DB.Exec(apkDescSql, apkId, v.PackageName, v.Title, v.Description, permission, v.PrivacyPolicy, v.RecentChanges, "ru", iconId, iconId, statusApkDesc)

		if err != nil {
			log.Printf("oz_apk_desc err:%v", err)
			continue
		}

		for _, val := range v.Screenshots {
			//写入oz_image表
			newScree := "http://18.177.149.123:8001/resources/" + v.PackageName + "/screenshots_"
			newScreeUrl := strings.Replace(val, "https://play-lh.googleusercontent.com/", newScree, 1)
			scrSql := "INSERT INTO kk_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
			screResult, err := DB.Exec(scrSql, v.PackageName+"_Screenshots", newScreeUrl, newScreeUrl, 50, "ru")

			if err != nil {
				log.Printf("oz_image 3 err:%v", err)
				continue
			}

			//写入oz_apk_image表
			imageId, _ := screResult.LastInsertId()
			ssSql := "INSERT INTO kk_apk_image (`apk_id`, `image_id`,`language`) VALUES (?,?,?);"
			_, err = DB.Exec(ssSql, apkId, imageId, "ru")

			if err != nil {
				log.Printf("oz_apk_image err:%v", err)
				continue
			}
		}
	}
	return nil
}
