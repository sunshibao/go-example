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
	AppType           string   `bson:"appType"`
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
	ShareUrl               string      `bson:"shareUrl"`
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
		fmt.Println("mysql????????????")
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
		v := s.Language["zh_CN"]
		appDetails := s.Language["zh_CN"].Detail.AppDetails
		downloadUrl := "https://newmarket3.tt286.com/180/apk/2021/08/02/8p9luaea9s/1627890026612.apk"
		statusApk := 1
		apkType := 0
		cateName := appDetails.CategoryName

		if appDetails.AppType == "GAME" {
			apkType = 1
			cateName = "GAME_" + cateName
		}
		//??????oz_apk???
		log.Printf("PackageName:%s-----num:%d\n", s.Package, k+skip)

		// ?????????????????????
		var apkId int64
		sql := "select apk_id from oz_apk where package_name = ?"
		DB.Get(&apkId, sql, s.Package)

		if apkId > 0 {
			apkSql := "update oz_apk set gp_down_url = ? where package_name = ?"
			DB.Exec(apkSql, v.ShareUrl, s.Package)

			err = insertDB(s, "be_BY", apkId) //?????????
			if err != nil {
				continue
			}
			err = insertDB(s, "uk_UA", apkId) //???????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "kk_KZ", apkId) //????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "ka_GE", apkId) //???????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "uz_UZ", apkId) //???????????????
			if err != nil {
				continue
			}

		} else {
			apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,company,file_size,download_num,apk_res_type,apk_type,status,gp_down_url)values(?,?,?,?,?,?,?,?,?,?,?,?)"
			appResult, err := DB.Exec(apkSql, s.Package, v.Title, appDetails.VersionCode, appDetails.VersionString, downloadUrl, appDetails.DeveloperName, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, cateName, apkType, statusApk, v.ShareUrl)

			if err != nil {
				log.Printf("oz_apk err:%v", err)
				return err
			}
			apkId, _ = appResult.LastInsertId()

			err = insertDB(s, "zh_CN", apkId) //????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "en_US", apkId) //??????
			if err != nil {
				continue
			}
			err = insertDB(s, "ru_RU", apkId) //??????
			if err != nil {
				continue
			}
			err = insertDB(s, "be_BY", apkId) //?????????
			if err != nil {
				continue
			}
			err = insertDB(s, "uk_UA", apkId) //???????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "kk_KZ", apkId) //????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "ka_GE", apkId) //???????????????
			if err != nil {
				continue
			}
			err = insertDB(s, "uz_UZ", apkId) //???????????????
			if err != nil {
				continue
			}

		}
	}
	return nil
}

// ????????????????????????
func insertDB(s MongoAppInfo, lang string, apkId int64) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	newIcon := "http://18.177.149.123:8001/new/" + s.Package

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

	cateName := appDetails.CategoryName
	if appDetails.AppType == "GAME" {
		cateName = "GAME_" + cateName
	}

	statusApkDesc := 1
	statusImg := 1

	v.Icon = strings.Replace(v.Icon, "https://play-lh.googleusercontent.com", newIcon, 1)

	//??????oz_image???
	iconSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	iconResult, err := DB.Exec(iconSql, s.Package+"_icon", v.Icon, v.Icon, 50, lang, statusImg)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}

	iconId, _ := iconResult.LastInsertId()

	//??????oz_image???
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
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,status,in_app_product,install_notes,apk_res_type)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, s.Package, v.Title, v.DescriptionHtml, permission, v.Annotations.PrivacyPolicyUrl, appDetails.RecentChangesHtml, lang, iconId, iconId, statusApkDesc, inAppProduct, installNotes, cateName)

	if err != nil {
		log.Printf("oz_apk_desc err:%v", err)
		return err
	}

	//????????????
	for _, val := range v.ImgList {
		val = strings.Replace(val, "https://play-lh.googleusercontent.com", newIcon, 1)
		//??????oz_image???
		scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
		screResult, err := DB.Exec(scrSql, s.Package+"_Screenshots", val, val, 50, lang)

		if err != nil {
			log.Printf("oz_image 3 err:%v", err)
			return err
		}
		//??????oz_apk_image???
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
