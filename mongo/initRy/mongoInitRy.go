package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type ImageInfo struct {
	ImageType int    `bson:"imageType"`
	ImageUrl  string `bson:"imageUrl"`
}

type Ds struct {
	Low int `json:"low"`
}

type Dc struct {
	Low int `json:"low"`
}

type AppInfo struct {
	DeveloperName     string   `bson:"developerName"`
	VersionCode       int      `bson:"versionCode"`
	VersionString     string   `bson:"versionString"`
	Permission        []string `bson:"permission"`
	PackageName       string   `bson:"packageName"`
	CategoryName      string   `bson:"categoryName"`
	RecentChangesHtml string   `bson:"recentChangesHtml"`
	InfoDownloadSize  Ds       `bson:"infoDownloadSize"`
	DownloadCount     Dc       `bson:"downloadCount"`
	InstallNotes      string   `bson:"installNotes"`
	InAppProduct      string   `bson:"inAppProduct"`
}

type Annotations struct {
	PrivacyPolicyUrl string `bson:"privacyPolicyUrl"`
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
	Annotations            Annotations `bson:"annotations"`
	DromotionalDescription string      `bson:"promotionalDescription"`
}

type MongoAppInfo struct {
	Package          string                 `bson:"package"`
	Language         map[string]LangAppInfo `bson:"language"`
	CheckVersionCode int                    `bson:"donwloadUrl"`
}

var DB *sqlx.DB

func GetDatabase() *sqlx.DB {
	return DB
}

var LocalCateName = []string{"策略", "动作", "赌场", "角色扮演", "教育", "街机", "竞速", "卡牌", "冒险", "模拟", "教育", "文字", "休闲", "益智", "音乐", "知识问答", "桌面和棋类"}

func main() {

	host := "18.177.149.123:27017"
	database := "package"
	mongodb, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 5000 {
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
		err = insertDB(s, "zh_CN")
		if err != nil {
			continue
		}
		err = insertDB(s, "en_US")
		if err != nil {
			continue
		}
		err = insertDB(s, "ru_RU")
		if err != nil {
			continue
		}
	}
	return nil
}

func insertDB(s MongoAppInfo, lang string) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	if lang == "zh_CN" {
		lang = "zh_cn"
	} else if lang == "en_US" {
		lang = "en"
	} else if lang == "ru_RU" {
		lang = "ru"
	} else {
		lang = "zh_cn"
	}

	statusApk := 1
	statusApkDesc := 1
	statusImg := 1

	downloadUrl := "https://newmarket3.tt286.com/180/apk/2021/08/02/8p9luaea9s/1627890026612.apk"

	apkType := 0
	if IsInSlice(LocalCateName, appDetails.CategoryName) {
		apkType = 1
	}

	apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,company,file_size,download_num,apk_type,status)values(?,?,?,?,?,?,?,?,?,?)"
	appResult, err := DB.Exec(apkSql, s.Package, v.Title, appDetails.VersionCode, appDetails.VersionString, downloadUrl, appDetails.DeveloperName, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, apkType, statusApk)

	if err != nil {
		log.Printf("oz_apk err:%v", err)
		return err
	}
	apkId, _ := appResult.LastInsertId()

	//写入oz_image表
	iconSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	iconResult, err := DB.Exec(iconSql, s.Package+"_icon", v.Icon, v.Icon, 50, lang, statusImg)

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

	inAppProduct := 0
	installNotes := 0
	if appDetails.InstallNotes != "" {
		installNotes = 1
	}
	if appDetails.InAppProduct != "" {
		inAppProduct = 1
	}

	permission := strings.Join(appDetails.Permission, ",")
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,status,in_app_product,install_notes)values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, s.Package, v.Title, v.DescriptionHtml, permission, v.Annotations.PrivacyPolicyUrl, appDetails.RecentChangesHtml, lang, iconId, iconId, statusApkDesc, inAppProduct, installNotes)

	if err != nil {
		log.Printf("oz_apk_desc err:%v", err)
		return err
	}

	for _, val := range v.ImgList {
		//写入oz_image表
		scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
		screResult, err := DB.Exec(scrSql, s.Package+"_Screenshots", val, val, 50, lang)

		if err != nil {
			log.Printf("oz_image 3 err:%v", err)
			return err
		}

		//写入oz_apk_image表
		imageId, _ := screResult.LastInsertId()
		ssSql := "INSERT INTO oz_apk_image (`apk_id`, `image_id`,`language`) VALUES (?,?,?);"
		_, err = DB.Exec(ssSql, apkId, imageId, lang)
		if err != nil {
			log.Printf("oz_apk_image err:%v", err)
			return err
		}

	}
	return nil
}

// Find获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
