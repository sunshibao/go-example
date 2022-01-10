package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"strings"

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
	limit := 10
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

//614ade639a408f2d02455321
func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {
	c := mongodb.DB(database).C("info_11_22")
	var mongoAppInfos []MongoAppInfo

	c.Find(nil).Sort("_id").Skip(skip).Limit(limit).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		return err
	}

	for k, s := range mongoAppInfos {
		//if s.Status == 0 {
		//	log.Printf("下线应用 PackageName:%s-----num:%d\n", s.Package, k+skip)
		//	continue
		//}
		v := s.Language["ru_RU"]
		appDetails := s.Language["ru_RU"].Detail.AppDetails
		statusApk := 1
		apkType := 0
		cateName := appDetails.CategoryName

		if appDetails.AppType == "GAME" {
			apkType = 1
			cateName = "GAME_" + cateName
		}
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", s.Package, k+skip)

		// 检查包是否存在
		var apkId int64
		sql := "select apk_id from oz_apk where package_name = ?"
		DB.Get(&apkId, sql, s.Package)

		if apkId > 0 {

		} else {
			apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,company,company_type,file_size,download_num,apk_res_type,apk_type,status,gp_down_url,age_limit)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			appResult, err := DB.Exec(apkSql, s.Package, v.Title, appDetails.VersionCode, appDetails.VersionString, v.ShareUrl, appDetails.DeveloperName, 1, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, cateName, apkType, statusApk, v.ShareUrl, v.Annotations.BadgeForLegacyRating.Major)

			if err != nil {
				log.Printf("oz_apk err:%v", err)
				continue
			} else {

				apkId, _ = appResult.LastInsertId()
				err = insertDB(s, "ru_RU", apkId) //俄语
				if err != nil {
					continue
				}
			}

		}
	}
	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	newIcon := "http://18.177.149.123:8001/pic/img_" + s.UpdateTime.Format("2006-01-02") + "/" + s.Package

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

	v.DescriptionHtml = TrimHtml(v.DescriptionHtml)
	appDetails.RecentChangesHtml = TrimHtml(appDetails.RecentChangesHtml)

	permission := strings.Join(appDetails.Permission, ",")
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,status,in_app_product,install_notes,apk_res_type)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, s.Package, v.Title, v.DescriptionHtml, permission, v.Annotations.PrivacyPolicyUrl, appDetails.RecentChangesHtml, lang, iconId, iconId, statusApkDesc, inAppProduct, installNotes, cateName)

	if err != nil {
		log.Printf("oz_apk_desc err:%v", err)
		return err
	}

	apkPerSql := "insert into oz_apk_permission (apk_id,permission)values(?,?)"
	_, err = DB.Exec(apkPerSql, apkId, permission)

	if err != nil {
		log.Printf("apkPerSql err:%v", err)
		return err
	}
	//图片去重，写入oz_image ,oz_apk_image
	//DelImage(apkId, newIcon, lang, s.Package, v.ImgList)
	//添加图片
	for _, val := range v.ImgList {
		val = strings.Replace(val, "https://play-lh.googleusercontent.com", newIcon, 1)
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

func DelImage(apkId int64, newIcon, lang, package2 string, imageList []string) {
	mm := make(map[string][]string) //去重
	for _, v := range imageList {
		if v == "" {
			continue
		}
		v = strings.Replace(v, "https://play-lh.googleusercontent.com", newIcon, 1)

		res, err := http.Get(v)
		if err != nil {
			fmt.Println("A error occurred!")
			continue
		}
		// defer后的为延时操作，通常用来释放相关变量
		defer res.Body.Close()

		pix, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}
		fileCode := bytesToHexString(pix)

		mm[fileCode] = append(mm[fileCode], v)

	}
	for _, vv := range mm {
		for kkk, vvv := range vv {
			if kkk == 0 {
				//写入oz_image表
				scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
				screResult, err := DB.Exec(scrSql, package2+"_Screenshots", vvv, vvv, 50, lang)

				if err != nil {
					log.Printf("oz_image 3 err:%v", err)
					continue
				}
				//写入oz_apk_image表
				imageId, _ := screResult.LastInsertId()
				ssSql := "INSERT INTO oz_apk_image (`apk_id`, `image_id`,`language`) VALUES (?,?,?);"
				_, err = DB.Exec(ssSql, apkId, imageId, lang)
				if err != nil {
					log.Printf("oz_apk_image err:%v", err)
					continue
				}
			}
		}
	}
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
